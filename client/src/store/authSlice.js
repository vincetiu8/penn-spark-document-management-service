import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import axios from 'axios'
import API_ROUTE from '../apiRoute'
import { history } from '../Routes'
import setAuthorizationToken from '../auth/authorization'
import { ErrNetworkErr } from './errors'
import { onError, toggleLoad } from './util'

const initialState = {
  loadStatus: false,
  isAuthenticated: false,
  userData: null
}

export const loginUser = createAsyncThunk(
  'auth/loginUser',
  async (loginDetails, { rejectWithValue }) => {
    try {
      const response = await axios.post(`${API_ROUTE}/login`, loginDetails)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    loadUserFromStorage: state => {
      const token = localStorage.getItem('token')
      const userData = localStorage.getItem('userData')
      if (token && userData) {
        setAuthorizationToken(token)
        state.userData = JSON.parse(userData)
        state.isAuthenticated = true
      }
      state.loadStatus = true
    },
    updateUserData: (state, action) => {
      state.userData = action.payload
      localStorage.setItem('userData', JSON.stringify(action.payload))
    },
    logoutUser: state => {
      state.isAuthenticated = false
      localStorage.removeItem('token')
      setAuthorizationToken(false)
      history.push('/login')
    }
  },
  extraReducers: {
    [loginUser.pending]: toggleLoad,
    [loginUser.fulfilled]: (state, action) => {
      state.loading = false
      state.isAuthenticated = true
      state.userData = action.payload.user_data
      localStorage.setItem('token', action.payload.token)
      localStorage.setItem('userData', JSON.stringify(action.payload.user_data))
      setAuthorizationToken(action.payload.token)
    },
    [loginUser.rejected]: onError
  }
})

export default authSlice.reducer
export const {
  logoutUser,
  loadUserFromStorage,
  updateUserData
} = authSlice.actions
