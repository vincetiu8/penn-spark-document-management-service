import { createAsyncThunk, createEntityAdapter, createSlice } from '@reduxjs/toolkit'
import axios from 'axios'
import API_ROUTE from '../apiRoute'
import { ErrNetworkErr } from './errors'
import fileDownload from 'js-file-download'
import { onError, toggleLoad } from './util'

const foldersAdapter = createEntityAdapter({
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

const initialState = foldersAdapter.getInitialState({
  function: '',
  updates: []
})

export const createFolder = createAsyncThunk(
  'folder/createFolder',
  async (folder, { rejectWithValue }) => {
    try {
      const response = await axios.post(`${API_ROUTE}/folders`, folder)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const getFolder = createAsyncThunk(
  'folder/getFolder',
  async (folderID, {
    dispatch,
    getState,
    rejectWithValue
  }) => {
    try {
      const response = await axios.get(`${API_ROUTE}/folders/${folderID}`)
      const state = getState()
      const parentID = response.data.parent_folder_id
      if (parentID !== 0 && !state.folders.ids.includes(parentID)) dispatch(getFolder(parentID))
      response.data.child_folders.forEach(folder => {
        if (!state.folders.ids.includes(folder.id)) dispatch(getFolder(folder.id))
      })
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const updateFolder = createAsyncThunk(
  'folder/updateFolder',
  async (folder, { rejectWithValue }) => {
    try {
      const response = await axios.put(`${API_ROUTE}/folders/${folder.id}`, folder)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const deleteFolder = createAsyncThunk(
  'folder/deleteFolder',
  async (ids, { rejectWithValue }) => {
    try {
      await axios.delete(`${API_ROUTE}/folders/${ids.id}`)
      return ids
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const createFile = createAsyncThunk(
  'file/createFile',
  async ({
    file,
    data
  }, { rejectWithValue }) => {
    try {
      const response = await axios.post(`${API_ROUTE}/files`, file)
      const formData = new FormData()
      formData.append('file', data)
      await axios.put(`${API_ROUTE}/file-data/${response.data.id}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const getFile = createAsyncThunk(
  'file/getFile',
  async (entity, { rejectWithValue }) => {
    try {
      const response = await axios({
        method: 'get',
        url: `${API_ROUTE}/file-data/${entity.id}`,
        responseType: 'blob'
      })
      fileDownload(response.data, entity.name)
      return null
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const updateFile = createAsyncThunk(
  'file/updateFile',
  async (file, { rejectWithValue }) => {
    try {
      console.log(file)
      const response = await axios.put(`${API_ROUTE}/files/${file.id}`, file)
      return response.data
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

export const deleteFile = createAsyncThunk(
  'file/deleteFile',
  async (ids, { rejectWithValue }) => {
    try {
      await axios.delete(`${API_ROUTE}/files/${ids.id}`)
      return ids
    } catch (err) {
      if (err.response === undefined) {
        return rejectWithValue(ErrNetworkErr)
      }
      return rejectWithValue(err.response.data.error)
    }
  }
)

const addUpdate = (state, entity) => {
  if (state.updates.length === 0) {
    state.updates.push(entity)
    return
  }

  const updatedAt = Date.parse(entity.updated_at)
  let l = 0
  let h = state.updates.length - 1
  let m = parseInt((l + h) / 2)
  while (l < h) {
    const otherUpdate = Date.parse(state.updates[m].updated_at)
    if (updatedAt === otherUpdate) break

    updatedAt > otherUpdate ? h = m - 1 : l = m + 1

    m = parseInt((l + h) / 2)
  }

  state.updates.splice(Date.parse(state.updates[m].updated_at) < updatedAt ? m : m + 1, 0, entity)
  for (let i = m; i < state.updates.length; i++) {
    if (state.updates[i].id === entity.id && state.updates.files === undefined === entity.files === undefined) {
      state.updates.splice(i, 1)
      break
    }
  }
}

const foldersSlice = createSlice({
  name: 'folders',
  initialState,
  reducers: {
    clearFolder: state => {
      foldersAdapter.removeAll(state)
      state.error = null
    },
    updateFunction: (state, action) => {
      state.function = action.payload
    },
    getFolder: getFolder
  },
  extraReducers: {
    [createFolder.pending]: toggleLoad,
    [createFolder.fulfilled]: (state, action) => {
      state.loading = false
      foldersAdapter.upsertOne(state, action.payload)
      state.entities[action.payload.parent_folder_id].child_folders.push(action.payload)
    },
    [createFolder.rejected]: onError,
    [getFolder.pending]: toggleLoad,
    [getFolder.fulfilled]: (state, action) => {
      state.loading = false
      const folder = action.payload
      foldersAdapter.upsertOne(state, folder)

      const yesterday = Date.now() - 86400000

      if (Date.parse(folder.updated_at) > yesterday) {
        addUpdate(state, folder)
      }

      folder.files.forEach(file => {
        if (Date.parse(file.updated_at) > yesterday) {
          addUpdate(state, file)
        }
      })
    },
    [getFolder.rejected]: onError,
    [updateFolder.pending]: toggleLoad,
    [updateFolder.fulfilled]: (state, action) => {
      state.loading = false
      const parentID = action.payload.parent_folder_id
      if (parentID === 0) {
        return
      }
      const parentFolder = state.entities[parentID]
      parentFolder.child_folders = parentFolder.child_folders.filter(folder => folder.id !== action.payload.id)
      parentFolder.child_folders.push(action.payload)

      const folder = state.entities[action.payload.id]
      if (folder === undefined) {
        return
      }
      action.payload.child_folders = folder.child_folders
      action.payload.files = folder.files
      foldersAdapter.updateOne(state, action.payload)
    },
    [updateFolder.rejected]: onError,
    [deleteFolder.pending]: toggleLoad,
    [deleteFolder.fulfilled]: (state, action) => {
      state.loading = false
      const parent = state.entities[action.payload.parentFolderID]
      parent.child_folders = parent.child_folders.filter(child => child.id !== action.payload.id)
      foldersAdapter.removeOne(state, action.payload.id)
    },
    [deleteFolder.rejected]: onError,
    [createFile.pending]: toggleLoad,
    [createFile.fulfilled]: (state, action) => {
      state.loading = false
      const folder = state.entities[action.payload.folder_id]
      folder.files.push(action.payload)
      foldersAdapter.updateOne(state, folder)
    },
    [createFile.rejected]: onError,
    [getFile.pending]: toggleLoad,
    [getFile.fulfilled]: state => {
      state.loading = false
    },
    [getFile.rejected]: onError,
    [updateFile.pending]: toggleLoad,
    [updateFile.fulfilled]: (state, action) => {
      state.loading = false
      const folderID = action.payload.folder_id
      const folder = state.entities[folderID]
      folder.files = folder.files.filter(file => file.id !== action.payload.id)
      folder.files.push(action.payload)
    },
    [updateFile.rejected]: onError,
    [deleteFile.pending]: toggleLoad,
    [deleteFile.fulfilled]: (state, action) => {
      state.loading = false
      const folder = state.entities[action.payload.folderID]
      folder.files = folder.files.filter(file => file.id !== action.payload.id)
    },
    [deleteFile.rejected]: onError
  }
})

export const {
  clearFolder,
  updateFunction
} = foldersSlice.actions
export default foldersSlice.reducer

export const { selectById: selectFolderById } = foldersAdapter.getSelectors(state => state.folders)
