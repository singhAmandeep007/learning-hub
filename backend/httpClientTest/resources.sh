#!/bin/bash

# Base URL and configuration
BASE_URL="http://localhost:8080/api/resources"
ADMIN_SECRET="your-admin-secret-key"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to make API call and check response
make_request() {
    local description=$1
    local curl_command=$2
    
    echo -e "\n${BLUE}Executing: $description${NC}"
    echo "Running: $curl_command"
    
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

# Function to seed resources
seed_resources() {
    echo -e "\n${YELLOW}Seeding resources...${NC}"
    
    # Create Article Resources
		make_request "Create Article Resource 1" "curl -X POST \"$BASE_URL\" \
				-H \"AdminSecret: $ADMIN_SECRET\" \
				-F \"title=Test Article 1\" \
				-F \"description=This is a test article 1\" \
				-F \"type=article\" \
				-F \"tags=test,article,blog\" \
				-F \"url=https://shorturl.at/llE4F\" \
				-F \"thumbnail=@./images/image1.webp\""

		make_request "Create Article Resource 2" "curl -X POST \"$BASE_URL\" \
				-H \"AdminSecret: $ADMIN_SECRET\" \
				-F \"title=Test Article 2\" \
				-F \"description=This is a test article 2\" \
				-F \"type=article\" \
				-F \"tags=test,article,blog,productivity,coding,programming\" \
				-F \"url=https://shorturl.at/jbcrY\" \
				-F \"thumbnail=@./images/image2.webp\""

		make_request "Create Article Resource 3 without thumbnail" "curl -X POST \"$BASE_URL\" \
				-H \"AdminSecret: $ADMIN_SECRET\" \
				-F \"title=Test Article 3\" \
				-F \"description=This is a test article 3\" \
				-F \"type=article\" \
				-F \"tags=test,article,blog,productivity,coding,programming,dev,golang\" \
				-F \"url=https://shorturl.at/jbcrY\""

		# Create PDF Resource
		make_request "Create PDF Resource 1" "curl -X POST \"$BASE_URL\" \
				-H \"AdminSecret: $ADMIN_SECRET\" \
				-F \"title=Test PDF Document 1\" \
				-F \"description=This is a test PDF document 1\" \
				-F \"type=pdf\" \
				-F \"tags=test,pdf,documentation\" \
				-F \"file=@./pdfs/pdf2.pdf\" \
				-F \"thumbnail=@./images/image2.webp\""

		make_request "Create PDF Resource 2 without thumbnail" "curl -X POST \"$BASE_URL\" \
				-H \"AdminSecret: $ADMIN_SECRET\" \
				-F \"title=Test PDF Document 2\" \
				-F \"description=This is a test PDF document 2\" \
				-F \"type=pdf\" \
				-F \"tags=test,pdf,documentation,creativity,test,preact,react\" \
				-F \"file=@./pdfs/pdf1.pdf\""

		# Create Video Resource
		make_request "Create Video Resource 1" "curl -X POST \"$BASE_URL\" \
				-H \"AdminSecret: $ADMIN_SECRET\" \
				-F \"title=Test Video Tutorial 1\" \
				-F \"description=This is a test video tutorial 1\" \
				-F \"type=video\" \
				-F \"tags=test,video,tutorial\" \
				-F \"file=@./videos/video1.mp4\" \
				-F \"thumbnail=@./images/image1.webp\""

		make_request "Create Video Resource 2 without thumbnail" "curl -X POST \"$BASE_URL\" \
			-H \"AdminSecret: $ADMIN_SECRET\" \
			-F \"title=Test Video Tutorial 2\" \
			-F \"description=This is a test video tutorial 2\" \
			-F \"type=video\" \
			-F \"tags=test,video,tutorial\" \
			-F \"file=@./videos/video1.mp4\""
}

# Function to update a resource
update_resource() {
    echo -e "\n${YELLOW}Update Resource${NC}"
    read -p "Enter resource ID to update: " resource_id
    read -p "Enter new title (press enter to skip): " new_title
    read -p "Enter new description (press enter to skip): " new_description
    read -p "Enter new tags (comma-separated, press enter to skip): " new_tags

    # Build the update command
    update_cmd="curl -X PATCH \"$BASE_URL/$resource_id\" -H \"AdminSecret: $ADMIN_SECRET\""
    
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
    
    make_request "Delete Resource" "curl -X DELETE \"$BASE_URL/$resource_id\" \
        -H \"AdminSecret: $ADMIN_SECRET\""
}

# Function to display menu
show_menu() {
    echo -e "\n${BLUE}=== Resource Management Menu ===${NC}"
    echo "1. Seed Resources"
    echo "2. Update Resource"
    echo "3. Delete Resource"
    echo "4. Exit"
    echo -e "${BLUE}=============================${NC}\n"
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
            echo -e "\n${GREEN}Exiting...${NC}"
            exit 0
            ;;
        *)
            echo -e "\n${RED}Invalid choice. Please try again.${NC}"
            ;;
    esac
done 