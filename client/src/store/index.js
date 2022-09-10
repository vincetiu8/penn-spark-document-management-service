import { configureStore } from '@reduxjs/toolkit'
import authReducer from './authSlice'
import foldersReducer from './foldersSlice'
import usersReducer from './usersSlice'
import userRolesReducer from './userRolesSlice'

export default configureStore({
  reducer: {
    auth: authReducer,
    folders: foldersReducer,
    users: usersReducer,
    userRoles: userRolesReducer
  }
})
