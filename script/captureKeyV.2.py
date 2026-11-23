import asyncio
from playwright.async_api import async_playwright
from seleniumbase import cdp_driver
import os

# This function works as a script so that the key can be rotated out when needed

KEY_FILE = r"C:\Users\jaide\Documents\Workspace\Projects\GoFindThatBrand\key.txt"
#KEY_FILE = "/Users/jaidendiaz/Documents/WorkSpace/GoFindThatBrand/key.txt"
async def main():

    driver = await cdp_driver.start_async(locale="en")
    endpointUrl = driver.get_endpoint_url()
    captured = {}
    async with async_playwright() as p:
        browser = await p.chromium.connect_over_cdp(endpointUrl)
        context = browser.contexts[0]
        page = context.pages[0]
        def handleRequest(request):
            if "algolia.net/1/indexes" in request.url and request.method == "POST":
                captured["apiKey"] = request.headers.get("x-algolia-api-key")

        page.on("request", handleRequest)
        await page.goto("https://www.grailed.com/categories/all")
        await page.mouse.wheel(0, 1500)
        await page.wait_for_timeout(3000)
        await page.mouse.wheel(0, 1500)
        await page.wait_for_timeout(3000)
        await page.mouse.wheel(0, 1500)
        await page.wait_for_timeout(3000)
        print(captured)
        driver.quit()
    tmp = KEY_FILE
    
    with open(tmp, "w") as f:
        f.write(captured["apiKey"])
        os.replace(tmp, KEY_FILE)

    

if __name__ == "__main__":
    asyncio.run(main())