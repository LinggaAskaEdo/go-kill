package scheduler

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/linggaaskaedo/go-kill/common/component/scheduler"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service/product"

	"github.com/openpcc/openpcc/uuidv7"
	"github.com/rs/zerolog"
)

var (
	adjectives = []string{
		"Pro", "Ultra", "Premium", "Elite", "Classic", "Modern", "Smart", "Advanced",
		"Deluxe", "Standard", "Professional", "Essential", "Superior", "Eco", "Digital",
		"Wireless", "Portable", "Compact", "Heavy-Duty", "Ergonomic", "Luxury", "Budget",
		"High-Performance", "Lightweight", "Durable", "Sleek", "Innovative", "Versatile",
	}

	productTypes = []string{
		"Mouse", "Keyboard", "Monitor", "Headphones", "Speaker", "Cable", "Adapter",
		"Charger", "Case", "Stand", "Hub", "Dock", "Camera", "Microphone", "Router",
		"Tablet", "Phone", "Laptop", "Watch", "Tracker", "Light", "Sensor", "Lock",
		"Thermostat", "Doorbell", "Controller", "Console", "Headset", "Chair", "Desk",
		"Lamp", "Fan", "Heater", "Blender", "Toaster", "Kettle", "Vacuum", "Scale",
		"Mattress", "Pillow", "Blanket", "Curtain", "Rug", "Painting", "Vase", "Clock",
		"Dumbbell", "Yoga Mat", "Treadmill", "Bike", "Tent", "Backpack", "Sleeping Bag",
		"Ball", "Racket", "Gloves", "Helmet", "Boots", "Jacket", "Pants", "Shirt",
		"Dress", "Skirt", "Sweater", "Socks", "Hat", "Scarf", "Belt", "Wallet",
		"Bag", "Sunglasses", "Ring", "Necklace", "Bracelet", "Earrings", "Book",
		"Magazine", "DVD", "Vinyl", "Game", "Puzzle", "Doll", "Car", "Train",
	}

	materials = []string{
		"Aluminum", "Steel", "Plastic", "Wood", "Glass", "Carbon Fiber",
		"Leather", "Fabric", "Titanium", "Copper", "Ceramic", "Bamboo",
	}

	colors = []string{
		"Black", "White", "Silver", "Gray", "Blue", "Red", "Green", "Yellow",
		"Pink", "Purple", "Orange", "Brown", "Beige", "Navy", "Teal",
	}
)

type ProductGeneratorJob struct {
	log            zerolog.Logger
	productService product.ProductServiceItf
	cfg            scheduler.Config
	rng            *rand.Rand
	categoryIDs    []string
	mu             sync.RWMutex
}

func NewProductGeneratorJob(log zerolog.Logger, productService product.ProductServiceItf, cfg scheduler.Config) *ProductGeneratorJob {
	return &ProductGeneratorJob{
		log:            log,
		productService: productService,
		cfg:            cfg,
		rng:            rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0)),
		categoryIDs:    make([]string, 0),
	}
}

func (j *ProductGeneratorJob) Name() string {
	return "product_generator_job"
}

func (j *ProductGeneratorJob) Schedule() string {
	return j.cfg.Cron
}

func (j *ProductGeneratorJob) Run(ctx context.Context) error {
	if !j.cfg.Enabled {
		zerolog.Ctx(ctx).Debug().Msg(j.cfg.Name + " is disabled")
		return nil
	}

	zerolog.Ctx(ctx).Info().Int("batch_size", j.cfg.BatchSize).Msg("Generating random products")

	successCount := 0
	for i := 0; i < j.cfg.BatchSize; i++ {
		dataProduct := j.generateProduct(ctx, i)
		qty, rsv := j.generateInventory(ctx)

		_, err := j.productService.CreateProduct(ctx, dataProduct, qty, rsv)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Product generation failed")
			return err
		}

		successCount++
		zerolog.Ctx(ctx).Debug().Str("name", dataProduct.Name).Msg("Product generation successfully")
	}

	zerolog.Ctx(ctx).Info().Int("success", successCount).Int("total", j.cfg.BatchSize).Msg("Product generation batch completed")

	return nil
}

func (j *ProductGeneratorJob) generateProduct(ctx context.Context, idx int) dto.CreateProductRequest {
	name := j.generateProductName()
	desc := j.generateDescription(name)
	price := j.generatePrice()
	sku := j.generateSKU(idx + 1)
	isActive := j.rng.Float64() > 0.1 // 90% active
	cats := j.generateCategories(ctx)

	zerolog.Ctx(ctx).Debug().Str("name", name).Str("desc", desc).Float64("price", price).Str("sku", sku).Bool("isActive", isActive).Strs("cats", cats).Send()

	return dto.CreateProductRequest{
		Name:        name,
		Description: desc,
		Price:       price,
		SKU:         sku,
		IsActive:    isActive,
		Categories:  cats,
	}
}

func (j *ProductGeneratorJob) generateProductName() string {
	parts := []string{
		adjectives[rand.IntN(len(adjectives))],
		productTypes[rand.IntN(len(productTypes))],
	}

	// 40% chance to add color
	if rand.Float64() > 0.6 {
		parts = append(parts, colors[rand.IntN(len(colors))])
	}

	// 30% chance to add material
	if rand.Float64() > 0.7 {
		parts = append(parts, materials[rand.IntN(len(materials))])
	}

	return strings.Join(parts, " ")
}

func (j *ProductGeneratorJob) generateDescription(name string) string {
	templates := []string{
		"High-quality %s with excellent features and durability. Perfect for everyday use.",
		"Premium %s designed for professionals. Outstanding performance and reliability.",
		"Innovative %s that combines style and functionality. Great value for money.",
		"Top-rated %s with advanced features. Ideal for both beginners and experts.",
		"Best-selling %s trusted by thousands. Exceptional quality and design.",
	}

	tmpl := templates[rand.IntN(len(templates))]

	return fmt.Sprintf(tmpl, strings.ToLower(name))
}

func (j *ProductGeneratorJob) generateSKU(index int) string {
	uid, _ := uuidv7.New()
	return fmt.Sprintf("SKU-%s-%06d", uid.String()[:8], index)
}

func (j *ProductGeneratorJob) generatePrice() float64 {
	// Price ranges: (min, max)
	ranges := [][]float64{
		{5, 20},     // Budget
		{20, 50},    // Mid-range
		{50, 100},   // Upper mid
		{100, 300},  // Premium
		{300, 1000}, // Luxury
	}

	r := ranges[rand.IntN(len(ranges))]
	price := r[0] + rand.Float64()*(r[1]-r[0])

	// Round to two decimal places
	return float64(int(price*100)) / 100
}

func (j *ProductGeneratorJob) generateInventory(ctx context.Context) (int, int) {
	var quantity, reserved int

	r := rand.Float64()
	switch {
	case r > 0.9: // 10% out of stock
		quantity = 0
		reserved = 0
	case r > 0.7: // 20% low stock
		quantity = rand.IntN(20) + 1
		reserved = rand.IntN(min(5, quantity) + 1)
	default: // 70% normal stock
		quantity = rand.IntN(481) + 20 // 20–500
		reserved = rand.IntN(min(50, quantity/5) + 1)
	}

	zerolog.Ctx(ctx).Debug().Int("quantity", quantity).Int("reserved", reserved).Send()

	return quantity, reserved
}

func (j *ProductGeneratorJob) generateCategories(ctx context.Context) []string {
	j.mu.Lock()
	defer j.mu.Unlock()

	if len(j.categoryIDs) > 0 {
		ids := make([]string, len(j.categoryIDs))
		copy(ids, j.categoryIDs)
		return j.selectRandomCategories(ids)
	}

	categories, err := j.productService.ListCategories(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("generate_categories: failed to list categories")
		return nil
	}

	zerolog.Ctx(ctx).Debug().Int("categories", len(categories)).Send()

	if len(categories) == 0 {
		j.categoryIDs = []string{}
		return nil
	}

	ids := make([]string, len(categories))
	for i, cat := range categories {
		ids[i] = cat.ID
	}

	j.categoryIDs = ids
	return j.selectRandomCategories(ids)
}

func (j *ProductGeneratorJob) selectRandomCategories(ids []string) []string {
	max := len(ids)
	if max == 0 {
		return nil
	}

	var n int
	if max == 1 {
		n = 1
	} else {
		n = j.rng.IntN(max-1) + 1
	}

	shuffled := make([]string, max)
	copy(shuffled, ids)
	j.rng.Shuffle(max, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:n]
}
