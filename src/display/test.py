import requests
from json import loads, JSONDecodeError

while True:
    resp = requests.get("http://127.0.0.1:8080/api/viewport/getImage")
    try:
        assert resp.status_code == 200
        image = loads(resp.text)["text"]
        print(image)
    except Exception as e:
        print(f"error {e!r}")
    input("next?")