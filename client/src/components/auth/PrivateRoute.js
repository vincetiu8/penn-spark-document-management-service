import React from 'react'
import { Redirect, Route } from 'react-router-dom'
import { useSelector } from 'react-redux'

export const PrivateRoute = ({
  component: Component,
  ...rest
}) => {
  const authState = useSelector(state => state.auth)

  return (
    <Route
      {...rest}
      render={props => {
        if (!authState.loadStatus) {
          return ''
        }
        return authState.isAuthenticated ? <Component {...props} /> : <Redirect to='/login'/>
      }}
    />
  )
}
