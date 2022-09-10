import { Route, Router, Switch } from 'react-router-dom'
import { createBrowserHistory } from 'history'
import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Toolbar } from '@material-ui/core'
import Login from './components/auth/Login'
import { loadUserFromStorage, logoutUser } from './store/authSlice'
import { Folder } from './components/folders/Folder'
import { NavBar } from './components/nav/NavBar'
import { ErrUserUnauthorized } from './store/errors'
import { clearFolder } from './store/foldersSlice'
import { clearUsers } from './store/usersSlice'
import { NavDrawer } from './components/nav/NavDrawer'
import { PrivateRoute } from './components/auth/PrivateRoute'
import { Users } from './components/users/Users'
import { UserRoles } from './components/userRoles/UserRoles'
import { clearUserRoles } from './store/userRolesSlice'
import { Dashboard } from './components/dashboard/Dashboard'

export const history = createBrowserHistory()

const Routes = () => {
  const authState = useSelector(state => state.auth)
  const foldersState = useSelector(state => state.folders)
  const usersState = useSelector(state => state.users)
  const userRolesState = useSelector(state => state.userRoles)
  const dispatch = useDispatch()

  const [navDrawerOpen, setNavDrawerOpen] = useState(false)

  const toggleNavDrawer = () => {
    if (!authState.isAuthenticated) return
    setNavDrawerOpen(!navDrawerOpen)
  }

  useEffect(() => {
    if (!authState.loadStatus) dispatch(loadUserFromStorage())
  }, [authState])

  useEffect(() => {
    if (foldersState.error === ErrUserUnauthorized || usersState.error === ErrUserUnauthorized || userRolesState.error === ErrUserUnauthorized) {
      dispatch(logoutUser())
      dispatch(clearFolder())
      dispatch(clearUsers())
      dispatch(clearUserRoles())
    }
  }, [foldersState, usersState, userRolesState])

  return (
    <Router history={history}>
      <div className='App'>
        <NavBar toggleNavDrawer={toggleNavDrawer}/>
        <NavDrawer open={navDrawerOpen} onClose={toggleNavDrawer}/>
        <main>
          <Toolbar/>
          <Switch>
            <Route exact path='/login' component={Login}/>

            <PrivateRoute exact path='/' component={Dashboard}/>

            <PrivateRoute exact path='/folders/:id' component={Folder}/>

            <PrivateRoute exact path='/users' component={Users}/>

            <PrivateRoute exact path='/user-roles' component={UserRoles}/>
          </Switch>
        </main>
      </div>
    </Router>
  )
}

export default Routes
