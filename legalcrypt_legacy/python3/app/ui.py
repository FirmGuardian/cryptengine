import random
from crypt import *
from tkinter import Tk, Label, Button, Entry, StringVar, DISABLED, NORMAL, END, W, E, filedialog
from PIL import Image, ImageTk

class EncryptorUI:
    def __init__(self, master):
        self.master = master
        master.title("FirmGuardian Encryptor")

        self.secret_number = random.randint(1, 100)
        self.guess = None
        self.num_guesses = 0
        self.logo = ImageTk.PhotoImage(Image.open('firmguardian.gif'))
        self.logo_panel = Label(self.master, image=self.logo)
        self.logo_panel.pack()


        self.message = "Government-grade encryption made easy!"
        self.info_message = ""

        self.label_text = StringVar()
        self.label_text.set(self.message)
        self.label = Label(master, textvariable=self.label_text, font="Helvetica 18", fg="#636466")

        self.info_label_text = StringVar()
        self.info_label_text.set(self.info_message)
        self.info_label = Label(master, textvariable=self.info_label_text, font = "Helvetica 12", fg="#636466")

        #vcmd = master.register(self.validate) # we have to wrap the command
        #v = StringVar(self.master, value='Please enter a password')

        #self.entry = Entry(master, validate="key", textvariable=v, validatecommand=(vcmd, '%P'))


        self.encrypt_button = Button(master, text="Encrypt", command=self.encrypt_file)
        self.decrypt_button = Button(master, text="Decrypt", command=self.decrypt_file)

        self.logo_panel.grid(row=0, column=0, columnspan=2, sticky=W+E,pady=(20,0), padx=(20,20))
        self.label.grid(row=1, column=0, columnspan=2, sticky=W+E,pady=(20,0), padx=(20,20))
        self.info_label.grid(row=2, column=0, columnspan=2, sticky=W+E,pady=(20,20), padx=(20,20))

        #self.entry.grid(row=3, column=0, columnspan=2, sticky=W+E, pady=(20,20), padx=(20,20))
        self.encrypt_button.grid(row=4, column=0, pady=(0,20), padx=(60,0))
        self.encrypt_button.config(width=10,font = "Helvetica 16", fg="#636466")
        self.decrypt_button.grid(row=4, column=1, pady=(0,20), padx=(0,60))
        self.decrypt_button.config(width=10, font = "Helvetica 16", fg="#636466")

    # not currently used, but may come in handy if we want password fields
    def reset(self):
        #self.entry.delete(0, END)
        #self.message = ""
        self.label_text.set(self.message)

        self.encrypt_button.configure(state=NORMAL)
        self.decrypt_button.configure(state=NORMAL)

    def encrypt_file(self):
        check_keypair_existence_and_create()
        self.info_message = filedialog.askopenfile(mode='r').name

        self.info_label_text.set(self.info_message + ' encrypted!')
        encrypt(self.info_message)

    def decrypt_file(self):
        res = check_keypair_existence()

        if(res):
            file_name = filedialog.askopenfile(mode='r').name
            decrypt(file_name)
            self.info_message = file_name.strip('_ENCRYPTED') + ' decrypted!'
            self.info_label_text.set(self.info_message)
        else:
            self.info_message='Error - Could not locate keys'
            self.info_label_text.set(self.info_message)

root = Tk()
root.resizable(width=False, height=False)
#root.configure(background='#ECF0F1')
my_gui = EncryptorUI(root)
root.mainloop()
