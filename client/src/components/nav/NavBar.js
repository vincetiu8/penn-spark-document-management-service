import { AppBar, IconButton, Toolbar, Typography } from '@material-ui/core'
import { Menu } from '@material-ui/icons'
import React from 'react'

export const NavBar = ({ toggleNavDrawer }) => {
  return (
    <div>
      <AppBar position='fixed' style={{ zIndex: 1251 }}>
        <Toolbar>
          <IconButton edge='start' color='inherit' onClick={toggleNavDrawer}>
            <Menu/>
          </IconButton>
          <Typography variant='h6' color='inherit'>
            Penn SPARK Document Management Service
          </Typography>
        </Toolbar>
      </AppBar>
    </div>
  )
}
