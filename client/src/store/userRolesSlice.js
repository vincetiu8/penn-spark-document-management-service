import { createAsyncThunk, createEntityAdapter, createSlice } from '@reduxjs/toolkit'
import axios from 'axios'
import API_ROUTE from '../apiRoute'
import { ErrNetworkErr } from './errors'
import { onError, toggleLoad } from './util'

const userRolesAdapter = createEntityAdapter({
  sortComparer: (a, b) => {
    if (a.id > b.id) {
      return -1
    }
    if (a.id < b.id) {
      return 1
    }
    return 0
  }
})

const initialState = userRolesAdapter.getInitialState({
  loadedData: false
})

export const createUserRole = createAsyncThunk(
  'userRole/createUser',
  async (userRole, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.post(`${API_ROUTE}/user-roles`, userRole)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const getUserRoles = createAsyncThunk(
  'userRole/getUserRoles',
  async (_, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.get(`${API_ROUTE}/user-roles`)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const updateUserRole = createAsyncThunk(
  'userRole/updateUserRole',
  async (userRole, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.put(`${API_ROUTE}/user-roles/${userRole.id}`, userRole)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const deleteUserRole = createAsyncThunk(
  'userRole/deleteUserRole',
  async (userRoleID, {
    rejectWithValue
  }) => {
    try {
      await axios.delete(`${API_ROUTE}/user-roles/${userRoleID}`)
      return userRoleID
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const addAccessRole = createAsyncThunk(
  'userRole/addAccessRole',
  async (accessRole, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.post(`${API_ROUTE}/access-roles`, accessRole)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const removeAccessRole = createAsyncThunk(
  'userRole/removeAccessRole',
  async (accessRole, {
    rejectWithValue
  }) => {
    try {
      await axios.delete(`${API_ROUTE}/access-roles/${accessRole.id}`)
      return accessRole
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

const userRolesSlice = createSlice({
  name: 'userRoles',
  initialState,
  reducers: {
    clearUserRoles: state => {
      userRolesAdapter.removeAll(state)
      state.error = null
      state.loadedData = false
    }
  },
  extraReducers: {
    [createUserRole.pending]: toggleLoad,
    [createUserRole.fulfilled]: (state, action) => {
      state.loading = false
      userRolesAdapter.upsertOne(state, action.payload)
    },
    [createUserRole.rejected]: onError,
    [getUserRoles.pending]: toggleLoad,
    [getUserRoles.fulfilled]: (state, action) => {
      state.loading = false
      state.loadedData = true
      userRolesAdapter.upsertMany(state, action.payload)
    },
    [getUserRoles.rejected]: onError,
    [updateUserRole.pending]: toggleLoad,
    [updateUserRole.fulfilled]: (state, action) => {
      state.loading = false
      userRolesAdapter.upsertOne(state, action.payload)
    },
    [updateUserRole.rejected]: onError,
    [deleteUserRole.pending]: toggleLoad,
    [deleteUserRole.fulfilled]: (state, action) => {
      state.loading = false
      userRolesAdapter.removeOne(state, action.payload)
    },
    [deleteUserRole.rejected]: onError,
    [addAccessRole.pending]: toggleLoad,
    [addAccessRole.fulfilled]: (state, action) => {
      state.loading = false
      const userRole = state.entities[action.payload.user_role_id]
      userRole.access_roles.push(action.payload)
    },
    [addAccessRole.rejected]: onError,
    [removeAccessRole.pending]: toggleLoad,
    [removeAccessRole.fulfilled]: (state, action) => {
      state.loading = false
      const userRole = state.entities[action.payload.user_role_id]
      userRole.access_roles = userRole.access_roles.filter(role => role.id !== action.payload.id)
    },
    [removeAccessRole.rejected]: onError
  }
})

export const { clearUserRoles } = userRolesSlice.actions
export default userRolesSlice.reducer

export const {
  selectById: selectUserRoleById
} = userRolesAdapter.getSelectors(state => state.userRoles)
