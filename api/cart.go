package api

import (
	"encoding/json"

	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saanvi-iyer/gobblego-backend/internal/cart"
	"github.com/saanvi-iyer/gobblego-backend/models"
	"gorm.io/gorm"
)

type CartHandler struct {
	Repo cart.Repository
	DB   *gorm.DB
}

type CartResponse struct {
	ItemID   uuid.UUID `json:"item_id"`
	ItemName string    `json:"item_name"`
	Image    string    `json:"image"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
	UserName string    `json:"user_name"`
	UserID   uuid.UUID `json:"user_id"`
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{
		DB:   db,
		Repo: cart.NewCartRepo(),
	}
}
func (h *CartHandler) GetAllCarts(c *fiber.Ctx) error {
	var carts []models.Cart
	if err := h.DB.Find(&carts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch carts"})
	}

	for i, cart := range carts {
		var items []models.CartItem
		if len(cart.Items) > 0 {
			if err := json.Unmarshal(cart.Items, &items); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to parse cart items"})
			}
		}

		// Collect all user IDs across all items
		var allUserIDs []uuid.UUID
		for _, item := range items {
			allUserIDs = append(allUserIDs, item.UserIDs...)
		}
		allUserIDs = uniqueUUIDs(allUserIDs)

		// Fetch all related users in a single query
		var users []models.User
		if len(allUserIDs) > 0 {
			if err := h.DB.Where("user_id IN ?", allUserIDs).Find(&users).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user details"})
			}
		}

		// Map userID → userName
		idToUsername := make(map[uuid.UUID]string)
		for _, u := range users {
			idToUsername[u.UserID] = u.UserName
		}

		var enrichedItems []fiber.Map
		for _, item := range items {
			var menuItem models.Menu
			if err := h.DB.Where("item_id = ?", item.ItemID).First(&menuItem).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch item details"})
			}

			// Map user IDs to names for this item
			var usernames []string
			for _, uid := range item.UserIDs {
				if name, ok := idToUsername[uid]; ok {
					usernames = append(usernames, name)
				}
			}

			enrichedItems = append(enrichedItems, fiber.Map{
				"item_id":   item.ItemID,
				"item_name": menuItem.ItemName,
				"image":     menuItem.Images,
				"price":     menuItem.Price,
				"quantity":  item.Quantity,
				"user_id":   item.UserIDs,
				"user_name": usernames,
			})
		}

		// Replace raw items with enriched version
		if newItems, err := json.Marshal(enrichedItems); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to marshal cart items"})
		} else {
			carts[i].Items = newItems
		}
	}

	return c.JSON(carts)
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	var req struct {
		ItemID   uuid.UUID `json:"item_id"`
		Quantity int       `json:"quantity"`
		UserID   uuid.UUID `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := h.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	if user.CartID == uuid.Nil {
		newCart := models.Cart{
			PaymentStatus: "pending",
			BillAmount:    0.0,
			Items:         []byte("[]"),
		}

		if err := h.DB.Create(&newCart).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create new cart"})
		}

		user.CartID = newCart.CartID
		if err := h.DB.Save(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update user with cart_id"})
		}
	}

	var cart models.Cart
	if err := h.DB.Where("cart_id = ?", user.CartID).First(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cart"})
	}

	var items []models.CartItem
	if len(cart.Items) > 0 {
		if err := json.Unmarshal(cart.Items, &items); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse cart items"})
		}
	}

	itemExists := false
	for i := range items {
		if items[i].ItemID == req.ItemID {
			// If user already exists in this item
			userAlreadyIn := false
			for _, uid := range items[i].UserIDs {
				if uid == req.UserID {
					userAlreadyIn = true
					break
				}
			}
			// Add user if not already in list
			if !userAlreadyIn {
				items[i].UserIDs = append(items[i].UserIDs, req.UserID)
			}
			// Update quantity
			items[i].Quantity += req.Quantity
			itemExists = true
			break
		}
	}

	if !itemExists {
		newItem := models.CartItem{
			ItemID:   req.ItemID,
			UserIDs:  []uuid.UUID{req.UserID},
			Quantity: req.Quantity,
		}
		items = append(items, newItem)
	}

	updatedItems, err := json.Marshal(items)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update cart items"})
	}

	cart.Items = updatedItems
	if err := h.DB.Save(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save updated cart"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Item added/updated in cart"})
}

func (h *CartHandler) UpdateCartItem(c *fiber.Ctx) error {
	var req struct {
		ItemID   uuid.UUID `json:"item_id"`
		Quantity int       `json:"quantity"`
		UserID   uuid.UUID `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := h.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	var cart models.Cart
	if err := h.DB.Where("cart_id = ?", user.CartID).First(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cart"})
	}

	var items []models.CartItem
	if len(cart.Items) > 0 {
		if err := json.Unmarshal(cart.Items, &items); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse cart items"})
		}
	}

	itemFound := false
	for i := range items {
		if items[i].ItemID == req.ItemID {
			itemFound = true

			alreadyPresent := false
			for _, uid := range items[i].UserIDs {
				if uid == req.UserID {
					alreadyPresent = true
					break
				}
			}
			if !alreadyPresent {
				items[i].UserIDs = append(items[i].UserIDs, req.UserID)
			}

			if req.Quantity > 0 {
				items[i].Quantity = req.Quantity
			} else {
				items = slices.Delete(items, i, i+1)
			}
			break
		}
	}

	// If item not found and quantity > 0, create new one with this user as first contributor
	if !itemFound && req.Quantity > 0 {
		newItem := models.CartItem{
			ItemID:   req.ItemID,
			UserIDs:  []uuid.UUID{req.UserID},
			Quantity: req.Quantity,
		}
		items = append(items, newItem)
	}

	// Save back to DB
	updatedItems, err := json.Marshal(items)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to marshal cart items"})
	}

	cart.Items = updatedItems
	if err := h.DB.Save(&cart).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save updated cart"})
	}

	// Resolve all user IDs to usernames
	var allUserIDs []uuid.UUID
	for _, item := range items {
		allUserIDs = append(allUserIDs, item.UserIDs...)
	}
	allUserIDs = uniqueUUIDs(allUserIDs)

	var users []models.User
	if err := h.DB.Where("user_id IN ?", allUserIDs).Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch usernames"})
	}

	idToUsername := make(map[uuid.UUID]string)
	for _, u := range users {
		idToUsername[u.UserID] = u.UserName
	}

	// Response
	type itemResponse struct {
		ItemID    uuid.UUID   `json:"item_id"`
		Quantity  int         `json:"quantity"`
		UserIDs   []uuid.UUID `json:"user_id"`
		UserNames []string    `json:"user_name"`
	}

	var res []itemResponse
	for _, item := range items {
		var usernames []string
		for _, uid := range item.UserIDs {
			if name, ok := idToUsername[uid]; ok {
				usernames = append(usernames, name)
			}
		}
		res = append(res, itemResponse{
			ItemID:    item.ItemID,
			Quantity:  item.Quantity,
			UserIDs:   item.UserIDs,
			UserNames: usernames,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":     "Cart item updated",
		"updatedCart": res,
	})
}

func uniqueUUIDs(input []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]bool)
	var result []uuid.UUID
	for _, id := range input {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}

func removeUUID(list []uuid.UUID, target uuid.UUID) []uuid.UUID {
	var result []uuid.UUID
	for _, id := range list {
		if id != target {
			result = append(result, id)
		}
	}
	return result
}
