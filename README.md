# Church Lower 3rds

A real-time Bible verse and speaker name display system for OBS (Open Broadcaster Software) designed for church services and live streaming.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)


<img width="1986" height="1288" alt="image" src="https://github.com/user-attachments/assets/e8ce40da-fd31-430f-9324-0c65e2b1f6d6" />


## Features

### 📖 Bible Verse Display
- **Complete KJV Bible** - All 66 books with 31,000+ verses
- **Smart Quick Input** - Type "Genesis 1:16" for instant verse selection
- **Real-time Autocomplete** - Intelligent book name suggestions with validation
- **Chapter Navigation** - Seamlessly browse through chapters and verses
- **Live Preview** - See verse content before displaying on stream
- **Next Verse Preview** - Preview upcoming verse for smooth transitions

### 🎤 Speaker Management
- **Speaker Names Display** - Professional lower third graphics for speaker identification
- **Autocomplete System** - Quick speaker selection with predefined names
- **Conflict Prevention** - Smart display logic prevents overlapping content
- **Customizable List** - Easy speaker management via text file

### 🎥 OBS Integration
- **Browser Source Ready** - Direct integration with OBS Studio
- **Fade Animations** - Smooth transitions for professional appearance
- **Dynamic Sizing** - Automatically adjusts for long verses
- **Bottom Positioning** - Non-intrusive placement for broadcast overlay
- **Real-time Updates** - Instant synchronization between control and display

### 🎨 User Interface
- **Dark/Light Themes** - Toggle between visual modes
- **Responsive Design** - Works on desktop and tablet devices
- **Keyboard Shortcuts** - Efficient navigation with arrow keys and Enter
- **Visual Feedback** - Color-coded status indicators and validation

## Quick Start

### Prerequisites
- [Go 1.21+](https://golang.org/dl/) installed on your system
- [OBS Studio](https://obsproject.com/) for broadcasting
- Modern web browser (Chrome, Firefox, Edge, Safari)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/RudeCoffee/ChurchLower3rds.git
   cd ChurchLower3rds
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start the server**
   ```bash
   go run main.go
   ```

4. **Access the interfaces**
   - **Control Interface**: http://localhost:8080/client.html
   - **OBS Browser Source**: http://localhost:8080/obs.html

### OBS Setup

1. **Add Browser Source**
   - In OBS, add a new "Browser" source
   - Set URL to: `http://localhost:8080/obs.html`
   - Recommended size: 1920x1080 (or your stream resolution)
   - Check "Shutdown source when not visible" for performance

2. **Position the Source**
   - Place the browser source above your main video feed
   - The lower thirds will appear at the bottom of the screen automatically

## Usage Guide

### Quick Verse Selection

The fastest way to display a Bible verse:

1. **Smart Input Method**
   ```
   Type: Genesis 1:16
   Press: Enter
   ```
   - Real-time validation prevents errors
   - Autocomplete suggests book names
   - Instant verse lookup and display

2. **Traditional Dropdown Method**
   - Select Book → Chapter → Verse
   - Use arrow keys for quick navigation
   - Preview appears automatically

### Speaker Display

1. **Add Speaker Names**
   - Edit `speakers.txt` file
   - Add one speaker name per line
   - Restart server to load new names

2. **Display Speaker**
   - Type speaker name (autocomplete enabled)
   - Click "Show Speaker" to display on stream
   - Click "Hide Speaker" to remove from display

### Navigation Controls

- **Previous/Next Verse**: Navigate through verses with buttons
- **Chapter Boundaries**: Automatically transitions between chapters
- **Keyboard Shortcuts**:
  - `Enter`: Select highlighted option
  - `Arrow Keys`: Navigate dropdown menus
  - `Escape`: Close dropdown menus

## File Structure

```
ChurchLower3rds/
├── main.go                 # Go server with WebSocket handling
├── client.html            # Control interface for operators
├── obs.html               # OBS browser source display
├── kjv.json               # Complete KJV Bible database (31,000+ verses)
├── speakers.txt           # Customizable speaker names list
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksums
└── README.md              # This documentation
```

## Configuration

### Speaker Names
Edit `speakers.txt` to customize available speakers:
```
Pastor John Smith
Elder Mary Johnson
Guest Speaker Mike Wilson
Worship Leader Sarah Davis
```

### Server Port
To change the default port (8080), modify the port in `main.go`:
```go
log.Println("Server starting on :8080")  // Change 8080 to your preferred port
```

### Display Styling
Customize the appearance by editing CSS in `obs.html`:
- Font sizes, colors, and positioning
- Background opacity and border radius
- Animation timing and effects

## Technical Details

### Architecture
- **Backend**: Go with Gorilla WebSocket for real-time communication
- **Frontend**: Vanilla HTML/CSS/JavaScript with WebSocket client
- **Data**: JSON-based Bible database with nested book/chapter/verse structure
- **Communication**: WebSocket messages for instant synchronization

### WebSocket Messages
The system uses these message types:
- `get_books`, `get_chapters`, `get_verses` - Data retrieval
- `preview_verse`, `get_verse` - Verse display commands
- `get_speakers` - Speaker list retrieval
- `show_speaker`, `hide_speaker` - Speaker display control

### Performance
- **Memory Usage**: ~50MB with full KJV Bible loaded
- **Response Time**: <10ms for verse lookups
- **Concurrent Users**: Supports multiple control interfaces
- **Browser Compatibility**: Modern browsers with WebSocket support

## Troubleshooting

### Common Issues

**Server won't start**
- Check if port 8080 is already in use
- Ensure Go is properly installed (`go version`)
- Verify all dependencies are installed (`go mod download`)

**OBS shows blank screen**
- Confirm the server is running
- Check the browser source URL: `http://localhost:8080/obs.html`
- Refresh the browser source in OBS
- Check browser console for JavaScript errors

**Verses not displaying**
- Verify `kjv.json` file is present and valid
- Check server console for loading errors
- Ensure WebSocket connection is established

**Autocomplete not working**
- Refresh the control interface
- Check network connectivity to server
- Verify JavaScript is enabled in browser

### Debug Mode
Add logging to troubleshoot issues:
```bash
go run main.go -verbose  # (if implemented)
```

## Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**
   - Follow Go best practices
   - Test thoroughly with OBS
   - Update documentation if needed
4. **Commit your changes**
   ```bash
   git commit -m "Add amazing feature"
   ```
5. **Push to your branch**
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Open a Pull Request**

### Development Guidelines
- Use `gofmt` for Go code formatting
- Test with multiple browsers for frontend changes
- Ensure WebSocket messages maintain backward compatibility
- Update README for new features

## License

This project is licensed under the MIT License - see the [LICENSES.md](LICENSES.md) file for details.

## Changelog

### Version 2.0.0 (Latest)
- ✨ Added smart quick verse input with autocomplete
- ✨ Real-time input validation and error prevention
- ✨ Dynamic background height adjustment for long verses
- ✨ Improved fade animations and visual feedback
- ✨ Enhanced keyboard navigation and shortcuts
- 🐛 Fixed chapter boundary navigation
- 🐛 Resolved speaker/verse display conflicts

### Version 1.0.0
- 📖 Complete KJV Bible integration
- 🎤 Speaker name display system
- 🎥 OBS browser source compatibility
- 🌓 Dark/light theme support
- ⚡ Real-time WebSocket communication

## Support

- **Issues**: [GitHub Issues](https://github.com/RudeCoffee/ChurchLower3rds/issues)
- **Documentation**: This README and inline code comments

## Acknowledgments

- King James Version Bible text (Public Domain)
- [Gorilla WebSocket](https://github.com/gorilla/websocket) for Go WebSocket implementation
- OBS Studio community for browser source capabilities

---

**Built with ❤️ for churches and ministries worldwide**
