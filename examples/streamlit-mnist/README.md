# Streamlit MNIST demo (drawable)

> The code is available [here](https://github.com/rahulsrma26/streamlit-mnist-drawable).

A simple digit recognition demo using [keras](https://www.tensorflow.org/overview) and [streamlit](https://www.streamlit.io/). It uses [streamlit-drawable-canvas](https://github.com/andfanilo/streamlit-drawable-canvas) for drawing on canvas.

![demo](img/demo.gif)

[streamlit](https://www.streamlit.io/) is an open-source app framework, which is the easiest way for data scientists and machine learning engineers to create beautiful, performant apps. All in pure Python, no longer fiddling with javascript.

This demo contains two parts: training a simple digit recognition model using mnist dataset and a webapp to live demo that model.
 
## Running demo

1. First install all the dependencies

    ```
    envd up
    ```

2. Train model

    Run all the cells of [train.ipynb](train.ipynb) manually.

3. Run demo web-app

    ```
    envd up -f build.envd:serve
    ```
