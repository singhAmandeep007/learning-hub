#!/bin/bash

PRODUCT="ecomm" # Generic e-commerce product namespace
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# Base URL and configuration
BASE_URL="http://localhost:8000/api/v1/${PRODUCT}/resources"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# --- Generic E-commerce Data Generation ---

ECOMMERCE_ACTIONS=("browse" "compare" "choose" "buy" "review" "save" "track" "return" "upgrade" "bundle")
ECOMMERCE_OBJECTS=("products" "orders" "categories" "carts" "wishlists" "coupons" "shipments" "payments" "reviews" "returns")
ECOMMERCE_CONCEPTS=("seasonal" "best-selling" "budget" "premium" "limited-time" "new-arrival" "bundle" "recommended")

ECOMMERCE_TAGS=(
    "electronics"
    "fashion"
    "home"
    "kitchen"
    "beauty"
    "fitness"
    "deals"
    "new arrivals"
    "best sellers"
    "gift ideas"
    "budget picks"
    "premium picks"
    "buying guide"
    "product comparison"
    "top rated"
    "customer favorites"
    "holiday sale"
    "back to school"
    "daily essentials"
    "smart home"
    "mobile accessories"
    "gaming"
    "office setup"
    "sustainable"
    "shipping"
    "returns"
    "warranty"
    "size guide"
    "care instructions"
    "how to choose"
    "unboxing"
    "installation"
    "maintenance"
    "troubleshooting"
    "faq"
    "secure checkout"
    "payment options"
    "order tracking"
    "customer support"
    "recommendations"
)

ARTICLE_TITLES=(
    "How to Choose the Right Smartphone for Everyday Use"
    "Laptop Buying Guide: Best Options by Budget"
    "Top 10 Home Appliances Worth Buying This Year"
    "Wireless Earbuds Comparison: Comfort, Battery, and Sound"
    "Best Smartwatch Features to Look for Before You Buy"
    "Gaming Setup Essentials: Monitors, Keyboards, and Mice"
    "How to Build a Productive Home Office on Any Budget"
    "TV Buying Guide: OLED vs QLED vs LED Explained"
    "Back-to-School Tech Checklist for Students"
    "Kitchen Essentials for First-Time Homeowners"
    "How to Pick the Perfect Gift for Any Occasion"
    "Camera Buying Guide for Beginners and Creators"
    "Top-Rated Fitness Gadgets for Home Workouts"
    "How to Compare Product Specs Like a Pro"
    "Care and Maintenance Tips to Extend Product Lifespan"
)

PDF_TITLES=(
    "E-commerce Product Catalog Overview"
    "Seasonal Deals and Promotions Handbook"
    "Customer Support and Return Policy Guide"
    "Warranty and Protection Plan Reference"
    "Product Care and Maintenance Manual"
    "Gift Buying Planner and Checklist"
    "Smart Home Starter Pack Guide"
    "Office Setup Planning Workbook"
    "Fashion Size and Fit Reference"
    "Secure Checkout and Payment Safety Guide"
)

VIDEO_TITLES=(
    "Top Budget Smartphones Compared in 10 Minutes"
    "Unboxing and First Impressions: Best-Selling Laptop"
    "How to Set Up a Smart Home Starter Kit"
    "Headphones Buying Tips: What Actually Matters"
    "Best Kitchen Gadgets for Everyday Cooking"
    "Gaming Accessories Explained for Beginners"
    "How to Pick the Right TV Size for Your Room"
    "Office Desk Setup: Productivity Essentials"
    "Quick Product Comparison: Premium vs Budget Picks"
    "How to Track Orders and Manage Returns Easily"
)

ARTICLE_DESCRIPTIONS=(
    "Practical buying guide covering key product features, value-for-money tips, and recommended options."
    "Step-by-step comparison advice to help shoppers evaluate specs, quality, and long-term usability."
    "Expert recommendations on choosing the right product category based on lifestyle, budget, and priorities."
    "Clear breakdown of must-have vs nice-to-have features so customers can make confident decisions."
    "Hands-on guidance for finding reliable products with strong ratings, warranty support, and durability."
)

PDF_DESCRIPTIONS=(
    "Comprehensive reference document with product information, policy details, and shopping best practices."
    "Detailed handbook for returns, warranties, and customer support workflows across common scenarios."
    "Concise catalog guide designed to help teams and customers navigate categories and offerings quickly."
    "Operational guide focused on secure payments, shipping expectations, and post-purchase support."
)

VIDEO_DESCRIPTIONS=(
    "Visual walkthrough that highlights product features, real-world usage, and practical buying recommendations."
    "Quick demo showing setup, comparison points, and tips to choose the best option for your budget."
    "Hands-on video session with use cases, do-and-don't guidance, and common troubleshooting advice."
    "Short tutorial focused on shopping decisions, product performance, and post-purchase care."
)

# Arrays of available local files to cycle through
IMAGE_FILES=(
    "$SCRIPT_DIR/images/ecommerce_sale_banner.gif"
    "$SCRIPT_DIR/images/ecommerce_product.jpeg"
    "$SCRIPT_DIR/images/ecommerce_product.jpg"
    "$SCRIPT_DIR/images/ecommerce_product.png"
    "$SCRIPT_DIR/images/ecommerce_product.webp"
)
PDF_FILES=("$SCRIPT_DIR/pdfs/ecommerce_catalog.pdf" "$SCRIPT_DIR/pdfs/ecommerce_simple.pdf")
VIDEO_FILES=("$SCRIPT_DIR/videos/ecommerce_promo.mp4" "$SCRIPT_DIR/videos/ecommerce_promo.webm")

# Function to get a random element from an array
get_random_element() {
    local arr=("$@")
    echo "${arr[$((RANDOM % ${#arr[@]}))]}"
}

# Function to get a realistic title based on content type
get_realistic_title() {
    local content_type=$1
    case $content_type in
        "article")
            get_random_element "${ARTICLE_TITLES[@]}"
            ;;
        "pdf")
            get_random_element "${PDF_TITLES[@]}"
            ;;
        "video")
            get_random_element "${VIDEO_TITLES[@]}"
            ;;
        *)
            echo "$(get_random_element "${ECOMMERCE_ACTIONS[@]}" | sed 's/.*/\u&/') $(get_random_element "${ECOMMERCE_OBJECTS[@]}"): E-commerce Best Practices"
            ;;
    esac
}

# Function to get a realistic description based on content type
get_realistic_description() {
    local content_type=$1
    case $content_type in
        "article")
            get_random_element "${ARTICLE_DESCRIPTIONS[@]}"
            ;;
        "pdf")
            get_random_element "${PDF_DESCRIPTIONS[@]}"
            ;;
        "video")
            get_random_element "${VIDEO_DESCRIPTIONS[@]}"
            ;;
        *)
            echo "Essential guide for e-commerce teams covering $(get_random_element "${ECOMMERCE_CONCEPTS[@]}") $(get_random_element "${ECOMMERCE_OBJECTS[@]}") best practices."
            ;;
    esac
}

# Function to generate realistic e-commerce-focused tags
get_realistic_ecommerce_tags() {
    local content_type=$1
    local num_tags=$((RANDOM % 4 + 4)) # Generate 4 to 7 tags
    local selected_tags=""
    
    # Always include content-type specific base tags
    case $content_type in
        "article")
            selected_tags="documentation,guide,tutorial"
            ;;
        "pdf")
            selected_tags="manual,reference,documentation"
            ;;
        "video")
            selected_tags="tutorial,demo,training"
            ;;
    esac
    
    # Create a copy of ECOMMERCE_TAGS to pick from
    local available_tags=("${ECOMMERCE_TAGS[@]}")
    
    # Add random e-commerce-specific tags
    for ((i=0; i<num_tags; i++)); do
        if [ ${#available_tags[@]} -eq 0 ]; then break; fi
        local random_index=$((RANDOM % ${#available_tags[@]}))
        selected_tags="$selected_tags,${available_tags[$random_index]}"
        # Remove the selected tag to avoid duplicates
        unset 'available_tags[random_index]'
        available_tags=("${available_tags[@]}") # Re-index the array
    done
    
    echo "$selected_tags"
}

# Function to make API call and check response
make_request() {
    local description=$1
    local curl_command=$2
    
    echo -e "\n${BLUE}Executing: $description${NC}"
    # echo "Running: $curl_command" # Optional: Uncomment for deep debugging
    
    response=$(eval "$curl_command")
    status=$?
    
    if [ $status -eq 0 ]; then
        echo -e "${GREEN}Success!${NC}"
        echo "Response: $response"
    else
        echo -e "${RED}Failed!${NC}"
        echo "Error: $response"
    fi
}

# Function to seed resources with realistic e-commerce data
seed_resources() {
    echo -e "\n${YELLOW}Seeding e-commerce resources...${NC}"
    
    read -p "Enter number of ARTICLE records to create: " num_articles
    read -p "Enter number of PDF records to create: " num_pdfs
    read -p "Enter number of VIDEO records to create: " num_videos

    # Create Article Resources with e-commerce content
    echo -e "\n${BLUE}Creating Article Resources...${NC}"
    for (( i=1; i<=num_articles; i++ )); do
        title=$(get_realistic_title "article")
        description=$(get_realistic_description "article")
        tags=$(get_realistic_ecommerce_tags "article")
        thumbnail=${IMAGE_FILES[$(( (i-1) % ${#IMAGE_FILES[@]} ))]}
        
        # Use realistic e-commerce content URLs
        urls=("https://example-store.com/blog/buying-guide" "https://example-store.com/blog/product-comparison" "https://example-store.com/blog/new-arrivals" "https://example-store.com/blog/shopping-tips")
        url=${urls[$(( (i-1) % ${#urls[@]} ))]}
        
        make_request "Create Article Resource $i: $(echo "$title" | cut -c1-50)..." "curl -s -X POST \"$BASE_URL\" \
            -F \"title=$title\" \
            -F \"description=$description\" \
            -F \"type=article\" \
            -F \"tags=$tags\" \
            -F \"url=$url\" \
            -F \"thumbnail=@$thumbnail\""
    done

    # Create PDF Resources with e-commerce content
    echo -e "\n${BLUE}Creating PDF Resources...${NC}"
    for (( i=1; i<=num_pdfs; i++ )); do
        title=$(get_realistic_title "pdf")
        description=$(get_realistic_description "pdf")
        tags=$(get_realistic_ecommerce_tags "pdf")
        file=${PDF_FILES[$(( (i-1) % ${#PDF_FILES[@]} ))]}
        thumbnail=${IMAGE_FILES[$(( (i-1) % ${#IMAGE_FILES[@]} ))]}

        make_request "Create PDF Resource $i: $(echo "$title" | cut -c1-50)..." "curl -s -X POST \"$BASE_URL\" \
            -F \"title=$title\" \
            -F \"description=$description\" \
            -F \"type=pdf\" \
            -F \"tags=$tags\" \
            -F \"file=@$file\" \
            -F \"thumbnail=@$thumbnail\""
    done

    # Create Video Resources with e-commerce content
    echo -e "\n${BLUE}Creating Video Resources...${NC}"
    for (( i=1; i<=num_videos; i++ )); do
        title=$(get_realistic_title "video")
        description=$(get_realistic_description "video")
        tags=$(get_realistic_ecommerce_tags "video")
        file=${VIDEO_FILES[$(( (i-1) % ${#VIDEO_FILES[@]} ))]}
        thumbnail=${IMAGE_FILES[$(( (i-1) % ${#IMAGE_FILES[@]} ))]}

        make_request "Create Video Resource $i: $(echo "$title" | cut -c1-50)..." "curl -s -X POST \"$BASE_URL\" \
            -F \"title=$title\" \
            -F \"description=$description\" \
            -F \"type=video\" \
            -F \"tags=$tags\" \
            -F \"file=@$file\" \
            -F \"thumbnail=@$thumbnail\""
    done
    
    echo -e "\n${GREEN}E-commerce resource seeding completed!${NC}"
    echo -e "${YELLOW}Created realistic e-commerce content with:${NC}"
    echo -e "  • $num_articles article resources"
    echo -e "  • $num_pdfs PDF guide resources" 
    echo -e "  • $num_videos video tutorial resources"
}

# Function to update a resource
update_resource() {
    echo -e "\n${YELLOW}Update Resource${NC}"
    read -p "Enter resource ID to update: " resource_id
    read -p "Enter new title (press enter to skip): " new_title
    read -p "Enter new description (press enter to skip): " new_description
    read -p "Enter new tags (comma-separated, press enter to skip): " new_tags

    # Build the update command
    update_cmd="curl -s -X PATCH \"$BASE_URL/$resource_id\""
    
    if [ ! -z "$new_title" ]; then
        update_cmd="$update_cmd -F \"title=$new_title\""
    fi
    if [ ! -z "$new_description" ]; then
        update_cmd="$update_cmd -F \"description=$new_description\""
    fi
    if [ ! -z "$new_tags" ]; then
        update_cmd="$update_cmd -F \"tags=$new_tags\""
    fi

    make_request "Update Resource" "$update_cmd"
}

# Function to delete a resource
delete_resource() {
    echo -e "\n${YELLOW}Delete Resource${NC}"
    read -p "Enter resource ID to delete: " resource_id
    
    make_request "Delete Resource" "curl -s -X DELETE \"$BASE_URL/$resource_id\""
}

# Function to display menu
show_menu() {
    echo -e "\n${BLUE}=== E-commerce Resource Management ===${NC}"
    echo "1. Seed E-commerce Resources (Realistic Data)"
    echo "2. Update Resource"
    echo "3. Delete Resource"
    echo "4. Exit"
    echo -e "${BLUE}=======================================${NC}\n"
}

# Main menu loop
while true; do
    show_menu
    read -p "Enter your choice (1-4): " choice

    case $choice in
        1)
            seed_resources
            ;;
        2)
            update_resource
            ;;
        3)
            delete_resource
            ;;
        4)
            echo -e "\n${GREEN}Exiting E-commerce Resource Seeder...${NC}"
            exit 0
            ;;
        *)
            echo -e "\n${RED}Invalid choice. Please try again.${NC}"
            ;;
    esac
done