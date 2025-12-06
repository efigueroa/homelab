#!/bin/bash

# System Information Gathering Script
# Outputs comprehensive system information to info.md for troubleshooting

OUTPUT_FILE="./info.md"

# Clear or create the output file
> "$OUTPUT_FILE"

echo "# System Information Report" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "Generated on: $(date)" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ==============================================================================
# OS Information
# ==============================================================================
echo "## Operating System Information" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Distribution info
if [ -f /etc/os-release ]; then
    echo "\`\`\`" >> "$OUTPUT_FILE"
    cat /etc/os-release >> "$OUTPUT_FILE"
    echo "\`\`\`" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
fi

# Kernel version
echo "**Kernel Version:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
uname -r >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Architecture
echo "**Architecture:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
uname -m >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ==============================================================================
# Desktop Environment Information
# ==============================================================================
echo "## Desktop Environment Information" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Current desktop session
echo "**Current Desktop Session:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "XDG_CURRENT_DESKTOP: ${XDG_CURRENT_DESKTOP:-Not set}" >> "$OUTPUT_FILE"
echo "XDG_SESSION_DESKTOP: ${XDG_SESSION_DESKTOP:-Not set}" >> "$OUTPUT_FILE"
echo "DESKTOP_SESSION: ${DESKTOP_SESSION:-Not set}" >> "$OUTPUT_FILE"
echo "GDMSESSION: ${GDMSESSION:-Not set}" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Installed desktop environments (check for common ones)
echo "**Installed Desktop Environments:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"

# Check for common desktop packages
for de in gnome-shell plasma-desktop xfce4-session mate-session cinnamon lxqt-session; do
    if command -v "$de" &> /dev/null || dpkg -l | grep -q "^ii.*$de" 2>/dev/null || rpm -q "$de" &> /dev/null; then
        echo "✓ $de found" >> "$OUTPUT_FILE"
    fi
done

echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ==============================================================================
# Display Server Information (Xorg vs Wayland)
# ==============================================================================
echo "## Display Server Information" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Currently running display server
echo "**Currently Running:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "XDG_SESSION_TYPE: ${XDG_SESSION_TYPE:-Not set}" >> "$OUTPUT_FILE"
echo "WAYLAND_DISPLAY: ${WAYLAND_DISPLAY:-Not set}" >> "$OUTPUT_FILE"
echo "DISPLAY: ${DISPLAY:-Not set}" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Check which is actually running via process
echo "**Running Display Server Processes:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
if pgrep -a Xorg > /dev/null 2>&1; then
    echo "Xorg is running:" >> "$OUTPUT_FILE"
    pgrep -a Xorg >> "$OUTPUT_FILE"
fi
if pgrep -a Xwayland > /dev/null 2>&1; then
    echo "Xwayland is running:" >> "$OUTPUT_FILE"
    pgrep -a Xwayland >> "$OUTPUT_FILE"
fi
if pgrep -a wayland > /dev/null 2>&1; then
    echo "Wayland compositor detected:" >> "$OUTPUT_FILE"
    pgrep -a wayland >> "$OUTPUT_FILE"
fi
# Check for common Wayland compositors
for compositor in gnome-shell kwin_wayland weston sway mutter; do
    if pgrep -a "$compositor" > /dev/null 2>&1; then
        echo "$compositor is running:" >> "$OUTPUT_FILE"
        pgrep -a "$compositor" >> "$OUTPUT_FILE"
    fi
done
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Xorg version
echo "**Xorg Version (if installed):**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
if command -v Xorg &> /dev/null; then
    Xorg -version 2>&1 | grep "X.Org" || echo "Xorg installed but version not detected" >> "$OUTPUT_FILE"
elif command -v X &> /dev/null; then
    X -version 2>&1 | grep "X.Org" || echo "X server installed but version not detected" >> "$OUTPUT_FILE"
else
    echo "Xorg not found in PATH" >> "$OUTPUT_FILE"
fi
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Wayland version/info
echo "**Wayland Libraries (if installed):**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
if ldconfig -p 2>/dev/null | grep -q "libwayland"; then
    ldconfig -p | grep libwayland >> "$OUTPUT_FILE"
else
    echo "Wayland libraries not found" >> "$OUTPUT_FILE"
fi
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ==============================================================================
# Steam Information
# ==============================================================================
echo "## Steam Information" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Check if Steam is installed
echo "**Steam Installation:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
if command -v steam &> /dev/null; then
    echo "Steam found in PATH: $(which steam)" >> "$OUTPUT_FILE"

    # Try to get Steam version
    if [ -f "$HOME/.steam/steam/package/steam_client_ubuntu_version.txt" ]; then
        echo "Steam client version: $(cat $HOME/.steam/steam/package/steam_client_ubuntu_version.txt)" >> "$OUTPUT_FILE"
    fi

    # Steam install locations
    echo "" >> "$OUTPUT_FILE"
    echo "Steam directories:" >> "$OUTPUT_FILE"
    [ -d "$HOME/.steam" ] && echo "  ~/.steam exists" >> "$OUTPUT_FILE"
    [ -d "$HOME/.local/share/Steam" ] && echo "  ~/.local/share/Steam exists" >> "$OUTPUT_FILE"

else
    echo "Steam not found in PATH" >> "$OUTPUT_FILE"

    # Check common installation paths
    if [ -d "$HOME/.steam" ] || [ -d "$HOME/.local/share/Steam" ]; then
        echo "Steam directories found but steam command not in PATH:" >> "$OUTPUT_FILE"
        [ -d "$HOME/.steam" ] && echo "  ~/.steam exists" >> "$OUTPUT_FILE"
        [ -d "$HOME/.local/share/Steam" ] && echo "  ~/.local/share/Steam exists" >> "$OUTPUT_FILE"
    fi
fi
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Steam runtime information
echo "**Steam Runtime:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
if [ -d "$HOME/.steam/steam/ubuntu12_32/steam-runtime" ]; then
    echo "Steam Runtime found at: ~/.steam/steam/ubuntu12_32/steam-runtime" >> "$OUTPUT_FILE"
    if [ -f "$HOME/.steam/steam/ubuntu12_32/steam-runtime/version.txt" ]; then
        echo "Runtime version: $(cat $HOME/.steam/steam/ubuntu12_32/steam-runtime/version.txt)" >> "$OUTPUT_FILE"
    fi
fi
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Steam launch options environment variables
echo "**Steam-related Environment Variables:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
env | grep -i steam || echo "No Steam-related environment variables found" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ==============================================================================
# Bashrc Analysis for Xorg/Wayland Configuration
# ==============================================================================
echo "## Bashrc Configuration Analysis" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "**Xorg/Wayland Related Configuration in ~/.bashrc:**" >> "$OUTPUT_FILE"
echo "\`\`\`bash" >> "$OUTPUT_FILE"

if [ -f "$HOME/.bashrc" ]; then
    # Look for Xorg/Wayland related settings
    grep -E -i "(xorg|wayland|display|xdg_session|x11|weston|sway|SDL_VIDEODRIVER|GDK_BACKEND|QT_QPA_PLATFORM|CLUTTER_BACKEND)" "$HOME/.bashrc" | grep -v "^#" || echo "No Xorg/Wayland related configuration found in ~/.bashrc" >> "$OUTPUT_FILE"
else
    echo "~/.bashrc not found" >> "$OUTPUT_FILE"
fi

echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# Check other shell config files too
echo "**Xorg/Wayland Related Configuration in Other Shell Files:**" >> "$OUTPUT_FILE"

for config_file in "$HOME/.bash_profile" "$HOME/.profile" "$HOME/.zshrc"; do
    if [ -f "$config_file" ]; then
        echo "" >> "$OUTPUT_FILE"
        echo "From $(basename $config_file):" >> "$OUTPUT_FILE"
        echo "\`\`\`bash" >> "$OUTPUT_FILE"
        grep -E -i "(xorg|wayland|display|xdg_session|x11|weston|sway|SDL_VIDEODRIVER|GDK_BACKEND|QT_QPA_PLATFORM|CLUTTER_BACKEND)" "$config_file" | grep -v "^#" || echo "No relevant configuration found" >> "$OUTPUT_FILE"
        echo "\`\`\`" >> "$OUTPUT_FILE"
    fi
done

echo "" >> "$OUTPUT_FILE"

# ==============================================================================
# Graphics Information
# ==============================================================================
echo "## Graphics Information" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "**GPU Information:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
if command -v lspci &> /dev/null; then
    lspci | grep -E "VGA|3D|Display" >> "$OUTPUT_FILE"
else
    echo "lspci not available" >> "$OUTPUT_FILE"
fi
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "**Graphics Driver Information:**" >> "$OUTPUT_FILE"
echo "\`\`\`" >> "$OUTPUT_FILE"
if command -v glxinfo &> /dev/null; then
    glxinfo | grep -E "OpenGL vendor|OpenGL renderer|OpenGL version" >> "$OUTPUT_FILE"
else
    echo "glxinfo not available (install mesa-utils to get this info)" >> "$OUTPUT_FILE"
fi
echo "\`\`\`" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ==============================================================================
# Completion Message
# ==============================================================================
echo "" >> "$OUTPUT_FILE"
echo "---" >> "$OUTPUT_FILE"
echo "*Report generation complete*" >> "$OUTPUT_FILE"

echo "System information has been written to: $OUTPUT_FILE"
