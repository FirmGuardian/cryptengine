import os
from Cryptodome.PublicKey import RSA
from Cryptodome.Random import get_random_bytes
from Cryptodome.Cipher import AES, PKCS1_OAEP

# GLOBAL VARIABLES
CUSTOMER_NAME = 'CUSTOMER_X'
KEYS_EXIST = False


def generate_key_pair():
    code = 'itsasecret'
    key = RSA.generate(2048)
    encrypted_key = key.exportKey(passphrase=code, pkcs=8, protection="scryptAndAES128-CBC")
    with open('/Users/alexbujduveanu/.ssh/' + CUSTOMER_NAME + '.bin', 'wb') as f:
        f.write(encrypted_key)
    with open('/Users/alexbujduveanu/.ssh/' + CUSTOMER_NAME + '_PUBKEY.bin', 'wb') as f:
        f.write(key.publickey().exportKey())
    KEYS_EXIST = True

def encrypt(in_file_name):
    with open(in_file_name + '_ENCRYPTED', 'wb') as out_file:
        recipient_key = RSA.import_key(open('/Users/alexbujduveanu/.ssh/' + CUSTOMER_NAME + '_PUBKEY.bin').read())
        session_key = get_random_bytes(16)

        cipher_rsa = PKCS1_OAEP.new(recipient_key)
        out_file.write(cipher_rsa.encrypt(session_key))

        cipher_aes = AES.new(session_key, AES.MODE_EAX)
        data_string = ''
        with open(in_file_name, 'rb') as input_file:
            data_bytes = input_file.read()

        data = data_bytes
        ciphertext, tag = cipher_aes.encrypt_and_digest(data)

        out_file.write(cipher_aes.nonce)
        out_file.write(tag)
        out_file.write(ciphertext)

        os.remove(in_file_name)


def decrypt(in_file_name):
    code = 'itsasecret'
    with open(in_file_name, 'rb') as fobj:
        private_key = RSA.import_key(open('/Users/alexbujduveanu/.ssh/' + CUSTOMER_NAME + '.bin').read(), passphrase=code)
        enc_session_key, nonce, tag, ciphertext = [fobj.read(x) for x in (private_key.size_in_bytes(), 16, 16, -1)]
        cipher_rsa = PKCS1_OAEP.new(private_key)
        session_key = cipher_rsa.decrypt(enc_session_key)
        cipher_aes = AES.new(session_key, AES.MODE_EAX, nonce)
        data = cipher_aes.decrypt_and_verify(ciphertext, tag)
        print(data)

        # write result to disk
        with open(in_file_name.strip('_ENCRYPTED'), 'wb') as out_file:
            out_file.write(data)

    os.remove(in_file_name)

def check_keypair_existence_and_create():
    if not os.path.isfile('/Users/alexbujduveanu/.ssh/' + CUSTOMER_NAME + '.bin'):
        # create the keys
        generate_key_pair()

def check_keypair_existence():
    if os.path.isfile('/Users/alexbujduveanu/.ssh/' + CUSTOMER_NAME + '.bin'):
        return True
    else:
        return False

'''
# this is being used as a module, so no main() is needed

def main():
    generate_key_pair()
    print("key pair done")
    encrypt()
    print("encrypt done")

    decrypt()
    print("decrypt done")

if __name__ == '__main__':
    main()
'''
