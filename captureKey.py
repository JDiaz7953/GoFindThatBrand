import os
print("SCRIPT LOADED")
KEY_FILE = r"C:\Users\jaide\Documents\Workspace\Projects\GoFindThatBrand\key.txt"

def request(flow):
    host = flow.request.pretty_host
    print("SEEN:", host)
    if "mnrwefss2q" in host:
        key = flow.request.headers.get("x-algolia-api-key")
        print("ALGOLIA REQUEST, key present:", bool(key))
        if key:
            try:
                tmp = KEY_FILE + ".tmp"
                with open(tmp, "w") as f:
                    f.write(key)
                os.replace(tmp, KEY_FILE)
                print("WROTE KEY FILE")
            except Exception as e:
                print("WRITE FAILED:", e)