import { createAsyncThunk, createEntityAdapter, createSlice } from '@reduxjs/toolkit'
import axios from 'axios'
import API_ROUTE from '../apiRoute'
import { ErrNetworkErr } from './errors'
import { onError, toggleLoad } from './util'

const usersAdapter = createEntityAdapter({
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

const initialState = usersAdapter.getInitialState({
  loadedData: false
})

export const createUser = createAsyncThunk(
  'user/createUser',
  async (user, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.post(`${API_ROUTE}/users`, user)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const getUsers = createAsyncThunk(
  'user/getUsers',
  async (_, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.get(`${API_ROUTE}/users`)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const getUser = createAsyncThunk(
  'user/getUser',
  async (userID, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.get(`${API_ROUTE}/users/${userID}`)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const updateUser = createAsyncThunk(
  'user/updateUser',
  async (user, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.put(`${API_ROUTE}/users/${user.id}`, user)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const deleteUser = createAsyncThunk(
  'user/deleteUser',
  async (userID, {
    rejectWithValue
  }) => {
    try {
      await axios.delete(`${API_ROUTE}/users/${userID}`)
      return userID
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const addUserRole = createAsyncThunk(
  'user/addUserRole',
  async ({
    userID,
    userRole
  }, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.post(`${API_ROUTE}/users/user-roles/${userID}`, userRole)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const removeUserRole = createAsyncThunk(
  'user/removeUserRole',
  async ({
    userID,
    userRole
  }, {
    rejectWithValue
  }) => {
    try {
      const response = await axios.delete(`${API_ROUTE}/users/user-roles/${userID}`, { data: userRole })
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

const usersSlice = createSlice({
  name: 'users',
  initialState,
  reducers: {
    clearUsers: state => {
      usersAdapter.removeAll(state)
      state.error = null
      state.loadedData = false
    }
  },
  extraReducers: {
    [createUser.pending]: toggleLoad,
    [createUser.fulfilled]: (state, action) => {
      state.loading = false
      usersAdapter.upsertOne(state, action.payload)
    },
    [createUser.rejected]: onError,
    [getUsers.pending]: toggleLoad,
    [getUsers.fulfilled]: (state, action) => {
      state.loading = false
      state.loadedData = true
      usersAdapter.upsertMany(state, action.payload)
    },
    [getUsers.rejected]: onError,
    [getUser.pending]: toggleLoad,
    [getUser.fulfilled]: (state, action) => {
      state.loading = false
      usersAdapter.upsertOne(state, action.payload)
    },
    [getUser.rejected]: onError,
    [updateUser.pending]: toggleLoad,
    [updateUser.fulfilled]: (state, action) => {
      state.loading = false
      usersAdapter.upsertOne(state, action.payload)
    },
    [updateUser.rejected]: onError,
    [deleteUser.pending]: toggleLoad,
    [deleteUser.fulfilled]: (state, action) => {
      state.loading = false
      usersAdapter.removeOne(state, action.payload)
    },
    [deleteUser.rejected]: onError,
    [addUserRole.pending]: toggleLoad,
    [addUserRole.fulfilled]: (state, action) => {
      state.loading = false
      usersAdapter.upsertOne(state, action.payload)
    },
    [addUserRole.rejected]: onError,
    [removeUserRole.pending]: toggleLoad,
    [removeUserRole.fulfilled]: (state, action) => {
      state.loading = false
      usersAdapter.upsertOne(state, action.payload)
    },
    [removeUserRole.rejected]: onError
  }
})

export const { clearUsers } = usersSlice.actions
export default usersSlice.reducer

export const {
  selectById: selectUserById
} = usersAdapter.getSelectors(state => state.users)
