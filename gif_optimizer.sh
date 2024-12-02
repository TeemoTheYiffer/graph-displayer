#!/bin/bash

# Check if gifsicle is installed
if ! command -v gifsicle &> /dev/null; then
    echo "gifsicle not found. Installing..."
    sudo apt-get update && sudo apt-get install -y gifsicle
fi

# Function to get file size in bytes
get_file_size() {
    stat -c%s "$1"
}

# Optimize all GIFs in the examples directory
for gif in ui/examples/*.gif; do
    if [[ -f "$gif" && "$gif" != *.optimized.gif ]]; then
        echo "Optimizing $gif..."
        gifsicle -O3 "$gif" > "${gif%.*}.optimized.gif"
        
        # Check if optimization was successful and smaller
        if [[ -f "${gif%.*}.optimized.gif" ]]; then
            original_size=$(get_file_size "$gif")
            optimized_size=$(get_file_size "${gif%.*}.optimized.gif")
            
            if (( optimized_size < original_size )); then
                mv "${gif%.*}.optimized.gif" "$gif"
                echo "Successfully optimized $gif ($original_size -> $optimized_size bytes)"
            else
                rm "${gif%.*}.optimized.gif"
                echo "Optimization didn't reduce size of $gif"
            fi
        fi
    fi
done

echo "GIF optimization complete!"
