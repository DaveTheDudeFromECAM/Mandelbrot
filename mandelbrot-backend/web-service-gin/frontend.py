import streamlit as st
import requests

# Add sliders to adjust the parameters of the Mandelbrot set.
xmin = st.slider("xmin", -2, 2, -2, 1)
xmax = st.slider("xmax", -2, 2, 2, 1)
ymin = st.slider("ymin", -2, 2, -2, 1)
ymax = st.slider("ymax", -2, 2, 2, 1)
iterations = st.slider("iterations", 50, 1000, 200, 1)

#@st.cache
def get_image_path():
    # Call the /mandelbrot endpoint to generate the PNG image.
    response = requests.get("http://localhost:8001/mandelbrot")
    st.write("it took ",response.json()["duration"], "to generate the set")
    return response.json()["imagePath"]



if st.button("Generate image"):
    # Get the file path of the generated PNG image.
    image_path = get_image_path()

    # Display the generated PNG image.
    st.image(image_path)
