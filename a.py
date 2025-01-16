from requests import post


r= post("https://gogramsess.vercel.app/api/gen", data={"appId": "123", "appHash": "456", "phoneNumber": "789"})

print(r.text)