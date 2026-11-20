import asyncio
from playwright.async_api import async_playwright
from seleniumbase import cdp_driver

async def main():
    driver = await cdp_driver.start_async(locale="en")
    endpointUrl = driver.get_endpoint_url()

    async with async_playwright() as p:
        browser = await p.chromium.connect_over_cdp(endpointUrl)
        page = browser.contexts[0].pages[0]
        await page.goto("https://www.grailed.com/")
        cookies = page.cookies()
        print(cookies)

if __name__ == "__main__":
    asyncio.run(main())