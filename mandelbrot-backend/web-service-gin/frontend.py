import streamlit as st
import requests

st.sidebar.title("Mandelbrot generation settings")

# Add sliders to adjust the parameters of the Mandelbrot set.
xmin = st.sidebar.slider("xmin", -2, 2, -2, 1)
xmax = st.sidebar.slider("xmax", -2, 2, 2, 1)
ymin = st.sidebar.slider("ymin", -2, 2, -2, 1)
ymax = st.sidebar.slider("ymax", -2, 2, 2, 1)
iterations = st.sidebar.slider("iterations", 50, 1000, 200, 1)

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
