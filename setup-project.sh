#!/bin/bash

SCRIPT_PATH=$(realpath "$0")
GITHUB_USERNAME="dbzer0"
OLD_NAME="go-rest-template"
CURRENT_DIR=$(basename "$PWD")

# Interactive input
echo "🚀 Go Project Setup"
echo "==================="

# Get the new project name (по умолчанию берём имя текущего каталога)
read -p "Enter new project name (default: $CURRENT_DIR): " NEW_NAME
if [ -z "$NEW_NAME" ]; then
    NEW_NAME="$CURRENT_DIR"
fi

echo -e "\n📋 Will replace both:"
echo "   $OLD_NAME → $NEW_NAME"
echo "   PROJECTNAME → $NEW_NAME"
echo "   GitHub path will be: github.com/$GITHUB_USERNAME/$NEW_NAME"
echo ""

read -p "Continue? (y/n): " CONFIRM
if [[ $CONFIRM != "y" && $CONFIRM != "Y" ]]; then
    echo "Setup cancelled."
    exit 0
fi

echo -e "\n🔄 Starting project setup..."

# 1. Replace go-rest-template with new name
echo "📝 Replacing technical repository name..."
find . -type f -not -path "*/\.git/*" -not -path "*/\.idea/*" -not -name "$(basename $SCRIPT_PATH)" | xargs grep -l "$OLD_NAME" 2>/dev/null | while read file; do
    echo "   Processing $file"
    sed -i '' "s|$OLD_NAME|$NEW_NAME|g" "$file"
done

# 2. Replace PROJECTNAME with new name
echo "📝 Replacing project name placeholder..."
find . -type f -not -path "*/\.git/*" -not -path "*/\.idea/*" -not -name "$(basename $SCRIPT_PATH)" | xargs grep -l "PROJECTNAME" 2>/dev/null | while read file; do
    echo "   Processing $file"
    sed -i '' "s|PROJECTNAME|$NEW_NAME|g" "$file"
done

# 3. Replace import paths specifically
echo "📦 Updating Go import paths..."
OLD_IMPORT_PATH="github.com/$GITHUB_USERNAME/$OLD_NAME"
NEW_IMPORT_PATH="github.com/$GITHUB_USERNAME/$NEW_NAME"

find . -type f -name "*.go" -not -path "*/\.git/*" -not -path "*/\.idea/*" | xargs grep -l "$OLD_IMPORT_PATH" 2>/dev/null | while read file; do
    echo "   Updating imports in $file"
    sed -i '' "s|$OLD_IMPORT_PATH|$NEW_IMPORT_PATH|g" "$file"
done

# 4. Remove .idea and .git directories
echo "🗑️  Removing .idea and .git directories..."
rm -rf .idea .git

# 5. Initialize a new git repository
echo "🔄 Initializing new git repository..."
git init

echo -e "\n✅ Project setup complete!"
echo "   - Project renamed to: $NEW_NAME"
echo "   - All occurrences of PROJECTNAME replaced"
echo "   - Import paths updated"
echo "   - .idea and .git directories removed"
echo "   - New git repository initialized"

# 6. Self-destruct
echo "🗑️  Removing setup script..."
rm "$SCRIPT_PATH"

echo -e "\n🚀 You're all set! Happy coding!"
