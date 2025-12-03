from playwright.sync_api import sync_playwright

def verify_frontend():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Navigate to the client page
        page.goto("http://localhost:8080/client.html")

        # Take a screenshot of the initial state
        page.screenshot(path="verification/client_initial.png")

        # Find the label that contains the AI switch input
        # The input is hidden, so we click the label or the slider
        # We look for the checkbox input with specific ID, then find its parent label
        ai_switch_label = page.locator("label.auto-advance-switch").filter(has=page.locator("#aiAssistantSwitch"))

        # Click the label to toggle
        ai_switch_label.click()

        # Check if the status text updated (Wait for it if needed)
        # The status update happens immediately in JS
        status = page.locator("#listeningStatus")

        # Wait for the text to appear to ensure JS ran
        page.wait_for_function("document.getElementById('listeningStatus').innerText.includes('Local AI is active')")

        print(f"Status text: {status.inner_text()}")

        # Take a screenshot after toggling
        page.screenshot(path="verification/client_ai_active.png")

        browser.close()

if __name__ == "__main__":
    verify_frontend()
