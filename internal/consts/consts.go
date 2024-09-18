package consts

// AppName stores the Application name
const (
	AppName          = "ecommerce"
	DatabaseType     = "postgres"
	AcceptedVersions = "v1"
)

// Context setting values
const (
	ContextAcceptedVersions       = "Accept-Version"
	ContextSystemAcceptedVersions = "System-Accept-Versions"
	ContextAcceptedVersionIndex   = "Accepted-Version-index"
	ContextErrorResponses         = "context-error-response"
	ContextLocallizationLanguage  = "lan"
)

// Pagination parameters
const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// Keys
const (
	Category     = "category"
	CategoryName = "category name"
	CategoryID   = "category id"

	Product     = "product"
	ProductName = "product name"
	ProductUrl  = "product image url"
	ProductID   = "product id"

	Variant              = "variant"
	VariantID            = "variant id"
	VariantMrp           = "mrp"
	VariantName          = "variant name"
	VariantDiscountPrice = "variant discout price"
	VariantQuantity      = "varaint quantity"

	Order         = "order"
	OrderID       = "order id"
	OrderItems    = "order items"
	OrderQuantity = "order quantity"
	OrderPrice    = "order price"
	OrderTotal    = "order total"
)

const (
	OrderStatusAccepted = "Accepted"
)

// Success Message
const (
	CategoryCreateSuccess = "Category created successfully"
	CategoryUpdateSuccess = "Category updated successfully"
	CategoryDeleteSuccess = "Category deleted successfully"
	CategoryFetchSuccess  = "Category fetched successfully"

	ProductCreateSuccess = "Product created successfully"
	ProductUpdateSuccess = "Product updated successfully"
	ProductDeleteSuccess = "Product deleted successfully"
	ProductFetchSuccess  = "Product fetched successfully"

	VariantCreateSuccess = "Variant created successfully"
	VariantUpdateSuccess = "Variant updated successfully"
	VariantDeleteSuccess = "Variant deleted successfully"
	VariantFetchSuccess  = "Variant fetched successfully"

	OrderCreateSuccess = "Order created successfully"
	OrderFetchSuccess  = "Order fetched successfully"
)

// Regex
const (
	UrlRegex = `^((ftp|http|https):\/\/)?([a-z0-9]+[.-])+[a-z0-9]{2,4}(:[0-9]{1,5})?(\/.*)?$`
)
