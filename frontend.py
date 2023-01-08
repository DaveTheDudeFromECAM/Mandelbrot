import streamlit as st
import requests

st.sidebar.title("Mandelbrot generation settings")

# Add sliders to adjust the parameters of the Mandelbrot set.
height = st.sidebar.slider("height", 500, 2000, 1000, 10)
width = st.sidebar.slider("width", 500, 2000, 1000, 10)
xmin = st.sidebar.slider("xmin", -2, 2, -2, 1)
xmax = st.sidebar.slider("xmax", -2, 2, 2, 1)
ymin = st.sidebar.slider("ymin", -2, 2, -2, 1)
ymax = st.sidebar.slider("ymax", -2, 2, 2, 1)
iterations = st.sidebar.slider("iterations", 50, 5000, 200, 1)

#@st.cache
def get_image():
    # Call the /mandelbrot endpoint to generate the PNG image.
    response = requests.get("http://localhost:8001/mandelbrot", params={
        "height":height,
        "width":width,
        "xmin": xmin,
        "xmax": xmax,
        "ymin": ymin,
        "ymax": ymax,
        "iterations": iterations
    })
    duration = response.json()["duration"]/1000000
    st.write("it took ",duration, "ms to generate the set")
    st.write({
        "height":height,
        "width":width,
        "xmin": xmin,
        "xmax": xmax,
        "ymin": ymin,
        "ymax": ymax,
        "iterations": iterations
    })
    return response.json()["imagePath"]

if st.button("Generate image"):
    # Get the file path of the generated PNG image.
    image_path = get_image()

    # Display the generated PNG image.
    st.image(image_path)
