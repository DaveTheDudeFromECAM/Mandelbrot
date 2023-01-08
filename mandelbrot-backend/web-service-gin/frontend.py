import streamlit as st
import requests

#@st.cache
def get_image_path():
    # Call the /mandelbrot endpoint to generate the PNG image.
    response = requests.get("http://localhost:8001/mandelbrot")
    return response.json()["imagePath"]

if st.button("Generate image"):
    # Get the file path of the generated PNG image.
    image_path = get_image_path()

    # Display the generated PNG image.
    st.image(image_path)
