
## 🛠️ Setup Instructions

Follow these simple steps to set up `kws_core`:

### 1️⃣ **Configuration**
First, configure your service by running:

```bash
make conf
```

This will generate the necessary configuration files for your cloud service. 📄✨

---

### 2️⃣ **Build**
Next, build the application with:

```bash
make build
```

This compiles all the source code and prepares `kws_core` for deployment. 🏗️🔧

---

### 3️⃣ **Start the Service**
Finally, start your service using:

```bash
sudo systemctl start kws_core
```

Your personal cloud is now live! 🎉🌈 You can check its status with:

```bash
sudo systemctl status kws_core
```

---

## 📂 Where Is Everything?

- **Configuration Files**: Located in the `conf/` directory.
- **Binary Files**: Generated in the `build/` directory after running `make build`.
- **Logs**: Check `/var/log/kws_core.log` for detailed logs.

---

## 💡 Tips & Tricks

- Want to restart the service? Run:
  ```bash
  sudo systemctl restart kws_core
  ```
- To stop the service gracefully, use:
  ```bash
  sudo systemctl stop kws_core
  ```

---

## 🎀 Why You'll Love It

- **Simple Setup**: Just three commands to get started!
- **Customizable**: Configure it to suit your needs.
- **Reliable**: Built to keep your personal cloud running smoothly.

---

Thank you for choosing `kws_core`! 🌟 Your cloud is now ready to serve you. Enjoy the journey! 🚀💕

--- 

Feel free to copy-paste this into an `README.md` file for your project! 😊

Sources

