import streamlit as st
import requests

# webpage content & instructions
st.header(  'This is a :red[_Mandelbrot_] set generation server  👨🏻‍💻')
st.caption( 'It uses _Golang_  for calculations and _python_ for the frontend with _Streamlit_ framework. ')
st.caption( '1. Open the sidebar menu using the arrow on the top left of the secreen \n'
            '2. Adjust the settings of the mandelbrot set and close the menu \n'
            '3. Click on "generate" and let the magic happen...  🪄✨ \n'
            )

# sliders to adjust the parameters of the Mandelbrot set
st.sidebar.title("Mandelbrot generation parameters")
height = st.sidebar.slider("image height", 500, 5000, 1000, 10)
width = st.sidebar.slider("image width", 500, 5000, 1000, 10)
xmin, xmax = st.sidebar.slider('X-axix plot limits',  2.0, 2.0, (-2.0, 1.5))
ymin, ymax = st.sidebar.slider('Y-axix plot limits', -2.0, 2.0, (-1.5, 1.5))
iterations = st.sidebar.slider("Number of iterations", 50, 5000, 200, 1)

def get_image():
    # request to backend
    response = requests.get("http://localhost:8001/mandelbrot", params={
        "iterations": iterations,
        "height":height,
        "width":width,
        "xmin": xmin,
        "xmax": xmax,
        "ymin": ymin,
        "ymax": ymax,
    })

    # generation time recieved form backend
    duration = response.json()["duration"]/1000000

    # mandelbrot params recap
    st.write(   ":t-rex: Oh ! It took ",duration, "ms to generate the set using folowing parameters: \n"
                "- Image resolution :",height, "pixels by",width,"pixels \n"
                "- X-axix plot range", xmin,"->", xmax, "\n"
                "- Y-axix plot range", ymin,"->", ymax, "\n"
                "- Iterations", iterations
    )
    return response.json()["imagePath"]

if st.button("Generate image"):
    # show spinner while waiting for the PNG generation
    with st.spinner(text='In progress'):

        # path of the generated PNG
        image_path = get_image()

        # display the PNG
        st.image(image_path)
        st.success(image_path)