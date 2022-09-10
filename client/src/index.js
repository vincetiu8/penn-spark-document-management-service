import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import Routes from './Routes'
import { BrowserRouter } from 'react-router-dom'
import * as serviceWorker from './serviceWorker'
import { Provider } from 'react-redux'
import store from './store/index'

ReactDOM.render((
  <BrowserRouter>
    <Provider store={store}>
      <Routes/>
    </Provider>
  </BrowserRouter>
), document.getElementById('root'))

serviceWorker.unregister()
