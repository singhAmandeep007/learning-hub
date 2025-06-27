#!/bin/bash

# Base URL and configuration
BASE_URL="http://localhost:8000/api/resources"
ADMIN_SECRET="your-admin-secret-key"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# --- Start of New Section: Random Data Generation ---

# Word and Tag lists for random data generation
WORDS=("Tech" "Future" "AI" "Code" "Productivity" "Design" "Creative" "Go" "React" "API" "System" "Life" "Hack" "Guide" "Tutorial")
TAG_LIST=("tech" "dev" "coding" "programming" "golang" "react" "javascript" "creative" "productivity" "blog" "guide" "testing" "api")

# Arrays of available local files to cycle through
IMAGE_FILES=("./images/image1.webp" "./images/image2.webp")
PDF_FILES=("./pdfs/pdf1.pdf" "./pdfs/pdf2.pdf")
VIDEO_FILES=("./videos/video1.mp4")

# Function to get a random element from an array
get_random_element() {
    local arr=("$@")
    echo "${arr[$((RANDOM % ${#arr[@]}))]}"
}

# Function to generate a random title
get_random_title() {
    echo "$(get_random_element "${WORDS[@]}") $(get_random_element "${WORDS[@]}"): A Deep Dive"
}

# Function to generate a random description
get_random_description() {
    echo "This is a detailed guide about $(get_random_element "${WORDS[@]}") and $(get_random_element "${WORDS[@]}") for modern developers."
}

# Function to generate a random comma-separated list of tags
get_random_tags() {
    local num_tags=$((RANDOM % 4 + 3)) # Generate 3 to 6 tags
    local tags_array=()
    # Ensure "test" tag is always present
    local selected_tags="test"
    
    # Create a copy of TAG_LIST to pick from
    local available_tags=("${TAG_LIST[@]}")

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

# --- End of New Section ---

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

# Function to seed resources (Transformed)
seed_resources() {
    echo -e "\n${YELLOW}Seeding resources...${NC}"
    
    read -p "Enter number of ARTICLE records to create: " num_articles
    read -p "Enter number of PDF records to create: " num_pdfs
    read -p "Enter number of VIDEO records to create: " num_videos

    # Create Article Resources
    for (( i=1; i<=num_articles; i++ )); do
        title=$(get_random_title)
        description=$(get_random_description)
        tags=$(get_random_tags)
        thumbnail=${IMAGE_FILES[$(( (i-1) % ${#IMAGE_FILES[@]} ))]}
        
        make_request "Create Article Resource $i" "curl -s -X POST \"$BASE_URL\" \
            -H \"AdminSecret: $ADMIN_SECRET\" \
            -F \"title=$title\" \
            -F \"description=$description\" \
            -F \"type=article\" \
            -F \"tags=$tags\" \
            -F \"url=https://shorturl.at/llE4F\" \
            -F \"thumbnail=@$thumbnail\""
    done

    # Create PDF Resources
    for (( i=1; i<=num_pdfs; i++ )); do
        title=$(get_random_title)
        description=$(get_random_description)
        tags=$(get_random_tags)
        file=${PDF_FILES[$(( (i-1) % ${#PDF_FILES[@]} ))]}
        thumbnail=${IMAGE_FILES[$(( (i-1) % ${#IMAGE_FILES[@]} ))]}

        make_request "Create PDF Resource $i" "curl -s -X POST \"$BASE_URL\" \
            -H \"AdminSecret: $ADMIN_SECRET\" \
            -F \"title=$title\" \
            -F \"description=$description\" \
            -F \"type=pdf\" \
            -F \"tags=$tags\" \
            -F \"file=@$file\" \
            -F \"thumbnail=@$thumbnail\""
    done

    # Create Video Resources
    for (( i=1; i<=num_videos; i++ )); do
        title=$(get_random_title)
        description=$(get_random_description)
        tags=$(get_random_tags)
        file=${VIDEO_FILES[$(( (i-1) % ${#VIDEO_FILES[@]} ))]}
        thumbnail=${IMAGE_FILES[$(( (i-1) % ${#IMAGE_FILES[@]} ))]}

        make_request "Create Video Resource $i" "curl -s -X POST \"$BASE_URL\" \
            -H \"AdminSecret: $ADMIN_SECRET\" \
            -F \"title=$title\" \
            -F \"description=$description\" \
            -F \"type=video\" \
            -F \"tags=$tags\" \
            -F \"file=@$file\" \
            -F \"thumbnail=@$thumbnail\""
    done
}

# Function to update a resource (Unchanged)
update_resource() {
    echo -e "\n${YELLOW}Update Resource${NC}"
    read -p "Enter resource ID to update: " resource_id
    read -p "Enter new title (press enter to skip): " new_title
    read -p "Enter new description (press enter to skip): " new_description
    read -p "Enter new tags (comma-separated, press enter to skip): " new_tags

    # Build the update command
    update_cmd="curl -s -X PATCH \"$BASE_URL/$resource_id\" -H \"AdminSecret: $ADMIN_SECRET\""
    
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

# Function to delete a resource (Unchanged)
delete_resource() {
    echo -e "\n${YELLOW}Delete Resource${NC}"
    read -p "Enter resource ID to delete: " resource_id
    
    make_request "Delete Resource" "curl -s -X DELETE \"$BASE_URL/$resource_id\" \
        -H \"AdminSecret: $ADMIN_SECRET\""
}

# Function to display menu (Unchanged)
show_menu() {
    echo -e "\n${BLUE}=== Resource Management Menu ===${NC}"
    echo "1. Seed Resources"
    echo "2. Update Resource"
    echo "3. Delete Resource"
    echo "4. Exit"
    echo -e "${BLUE}=============================${NC}\n"
}

# Main menu loop (Unchanged)
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
            echo -e "\n${GREEN}Exiting...${NC}"
            exit 0
            ;;
        *)
            echo -e "\n${RED}Invalid choice. Please try again.${NC}"
            ;;
    esac
done