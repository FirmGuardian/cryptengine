# LegalCrypt
Encryption solution for small to medium sized law firms. 

## Project Layout
The app is located in `/legalcrypt/python3/app/ui.py` and `/legalcrypt/python3/app/crypt.py`. 

The `/legalcrypt/python3/networking/` folder contains a PoC for a multithreaded HTTPS server which can accept up to 5 concurrent connections. `/legalcrypt/python3/networking/client.py` is the means through which `/legalcrypt/python3/app/ui.py` will connect to our server to do authentication and validation.  

### To run the app:
`python3 ui.py`

Yes, there will be bugs. 
They won't be pretty. 


## Project To Do List
- Implement file encryption using the programmatically generated RSA keys. 
  - Right now for Python3 the only encryption type supported is for text files.
  - Python2 handles data types differently so this was much easier on the earlier version. Example included under `/legalcrypt/python2/`
- Add a server component to listen for connections (in progress under the `/legalcrypt/python3/networking/` folder)
- Implement an authentication scheme for the client/server components
- Add context for alerts in the app (ex: green for success, red for failure)
- Figure out how to do storage of credentials and how the passwords will be hashed/stored (my vote is something like MySQL)
- Find out how to re-skin the desktop buttons and app for new logo/color scheme


## Dependencies
This encryption applet is dependent upon several Python packages. You can install them with the commands shown below. The `pycryptodome` package is actively maintained and allows an easy to use interface for RSA public-key cryptography. The `pillow` package is for image processing which is compatible with tkinter (the built-in Python GUI library). 

Use `pip3` to install the dependencies:
```
  pip3 install pycryptodome
  pip3 install pillow
```

> Note: check your default system Python version. On OS X it is set to 2.7. We decided to use 3.x to "future-proof" our application 
