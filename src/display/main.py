from json import loads
import tkinter as tk
from tkinter import filedialog
from PIL import Image, ImageTk, ImageDraw
import requests

def interpolate_color(t: float) -> tuple[int, int, int]:
    t = max(0.0, min(1.0, t))
    
    if t < 0.33:
        r = int(255 * (1 - t/0.33))
        g = int(255 * (t/0.33))
        b = 0
    elif t < 0.66:
        r = 0
        g = int(255 * (1 - (t-0.33)/0.33))
        b = int(255 * ((t-0.33)/0.33))
    else:
        r = int(255 * ((t-0.66)/0.34))
        g = 0
        b = int(255 * (1 - (t-0.66)/0.34))
    
    return (r, g, b)

def char_2_color(char: str) -> tuple[int, int, int]:
    charord = ord(char)
    if char == "#":
        return (0,0,0)
    elif char == "*":
        return (255, 255, 0)
    elif char == " ":
        return (255,255,255)
    if 65 <= charord <= 90:
        t = (charord - 65) / 26
        return interpolate_color(t)
    if 97 <= charord <= 122:
        t = (charord - 97) / 26
        color = interpolate_color(t)
        resize = 2.5
        return (
            int(color[0] / resize),
            int(color[1] / resize),
            int(color[2] / resize)
        )
    
    return (100,100,100)

class ViewPortImage:
    def __init__(self, from_string: str):
        self.from_string = from_string[:-1]
        self.splitted = self.from_string.split("\n")[1:]
        self.width = len(self.splitted[0])
        self.height = 0
        self.content: dict[tuple[int, int], 
                           tuple[int, int, int]] = {}
        for y, row in enumerate(self.splitted):
            for x, char in enumerate(row):
                self.content[x, y] = char_2_color(char)
            if len(row) == self.width:
                self.height += 1
                

class ImageEditor:
    def __init__(self, root):
        self.root = root
        
        # Создаем Canvas для изображения
        self.canvas = tk.Canvas(root, width=600, height=400, bg='white')
        self.canvas.pack()
        
        self.btn_frame = tk.Frame(root)
        self.btn_frame.pack(pady=10)
        
        self.load_btn = tk.Button(self.btn_frame, text="Пауза", 
                                  command=self.toggleGameLoopRunning)
        self.load_btn.pack(side=tk.LEFT, padx=5)
        
        # Переменные для хранения изображения
        self.image = None
        self.tk_image = None
        self.draw = None
        
        # Запускаем периодическую функцию
        self.periodic_task()
    
    def periodic_task(self):
        try:
            resp = requests.get("http://127.0.0.1:8080/api/viewport/getImage")
            image_string = loads(resp.text)["text"]
            
            try:
                viewport_image = ViewPortImage(image_string)
            except Exception as e:
                print(f"Ошибка в viewport_image {e!r}")
                raise
            
            try:
                self.load_image(viewport_image)
            except Exception as e:
                print(f"Ошибка в load_image {e!r}")
                raise
            
        except Exception as e:
            print(f"error {e!r}")

        self.root.after(50, self.periodic_task)
    
    def load_image(self, image_data: ViewPortImage, scale=12):
        """Загружает изображение из массива байтов [width, height, r0, g0, b0, ...]"""
        width, height = image_data.width, image_data.height
        
        img = Image.new("RGB", (width, height))
        pixels: dict = img.load()  # type: ignore


        for y in range(height):
            for x in range(width):
                pixels[x, y] = image_data.content[x, y]
        
        new_width, new_height = width, height
        if scale != 1.0:
            new_width = int(width * scale)
            new_height = int(height * scale)
            img = img.resize((new_width, new_height), Image.Resampling.NEAREST)
        
        # Сохраняем изображение и отображаем его
        self.image = img
        self.tk_image = ImageTk.PhotoImage(self.image)
        self.canvas.create_image(0, 0, anchor="nw", image=self.tk_image)
        self.canvas.config(width=new_width, height=new_height)
    
    def toggleGameLoopRunning(self):
        requests.post("http://127.0.0.1:8080/api/counter/toggleGameLoopRunning")


        
if __name__ == "__main__":
    root = tk.Tk()
    app = ImageEditor(root)
    root.mainloop()