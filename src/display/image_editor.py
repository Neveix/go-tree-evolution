import tkinter as tk
from tkinter import filedialog
from PIL import Image, ImageTk, ImageDraw

class ImageEditor:
    def __init__(self, root):
        self.root = root
        self.root.title("Редактор изображений")
        
        # Создаем Canvas для изображения
        self.canvas = tk.Canvas(root, width=600, height=400, bg='white')
        self.canvas.pack()
        
        # Переменные для хранения изображения
        self.image = None
        self.tk_image = None
        self.draw = None
        
        # Создаем кнопки
        self.btn_frame = tk.Frame(root)
        self.btn_frame.pack(pady=10)
        
        self.load_btn = tk.Button(self.btn_frame, text="Загрузить", command=self.load_image)
        self.load_btn.pack(side=tk.LEFT, padx=5)
        
        self.draw_btn = tk.Button(self.btn_frame, text="Рисовать", command=self.start_drawing)
        self.draw_btn.pack(side=tk.LEFT, padx=5)
        
        self.clear_btn = tk.Button(self.btn_frame, text="Очистить", command=self.clear_canvas)
        self.clear_btn.pack(side=tk.LEFT, padx=5)
        
        # Переменные для рисования
        self.drawing = False
        self.last_x = None
        self.last_y = None
        
    def load_image(self):
        file_path = filedialog.askopenfilename()
        if file_path:
            self.image = Image.open(file_path)
            self.image = self.image.resize((600, 400), Image.LANCZOS) # type: ignore
            self.tk_image = ImageTk.PhotoImage(self.image)
            self.canvas.create_image(0, 0, anchor=tk.NW, image=self.tk_image)
            self.draw = ImageDraw.Draw(self.image)
    
    def start_drawing(self):
        self.drawing = not self.drawing
        if self.drawing:
            self.draw_btn.config(text="Стоп")
            self.canvas.bind("<B1-Motion>", self.paint)
            self.canvas.bind("<ButtonRelease-1>", self.reset)
        else:
            self.draw_btn.config(text="Рисовать")
            self.canvas.unbind("<B1-Motion>")
            self.canvas.unbind("<ButtonRelease-1>")
    
    def paint(self, event):
        if self.drawing and self.draw:
            x, y = event.x, event.y
            if self.last_x and self.last_y:
                # Рисуем на изображении
                self.draw.line([self.last_x, self.last_y, x, y], fill='red', width=3)
                # Рисуем на холсте
                self.canvas.create_line(self.last_x, self.last_y, x, y, fill='red', width=3)
            self.last_x = x
            self.last_y = y
    
    def reset(self, event):
        self.last_x = None
        self.last_y = None
    
    def clear_canvas(self):
        self.canvas.delete("all")
        if self.image:
            self.image = Image.new("RGB", (600, 400), "white")
            self.tk_image = ImageTk.PhotoImage(self.image)
            self.canvas.create_image(0, 0, anchor=tk.NW, image=self.tk_image)
            self.draw = ImageDraw.Draw(self.image)

if __name__ == "__main__":
    root = tk.Tk()
    app = ImageEditor(root)
    root.mainloop()