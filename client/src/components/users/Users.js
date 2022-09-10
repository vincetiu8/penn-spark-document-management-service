import { useDispatch, useSelector } from 'react-redux'
import React, { useEffect, useState } from 'react'
import { getUsers } from '../../store/usersSlice'
import {
  FormControl,
  InputLabel,
  LinearProgress,
  makeStyles,
  MenuItem,
  Select,
  Table,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField
} from '@material-ui/core'
import { CreateUserBar } from './CreateUserBar'
import { EditUserDialog } from './EditUserDialog'
import { UserInfo } from './UserInfo'
import { ChangePasswordDialog } from './ChangePasswordDialog'
import { DeleteUserDialog } from './DeleteUserDialog'
import { AddUserRoleDialog } from './AddUserRoleDialog'

const columns = [
  {
    field: 'username',
    headerName: 'username',
    width: 150
  },
  {
    field: 'first_name',
    headerName: 'first name',
    width: 150
  },
  {
    field: 'last_name',
    headerName: 'last name',
    width: 150
  },
  {
    field: 'is_admin',
    headerName: 'admin',
    width: 150
  }
]

const useStyles = makeStyles({
  div: {
    width: '100%',
    height: '100%',
    position: 'fixed'
  },
  select: {
    minWidth: 150
  }
})

export const Users = () => {
  const dispatch = useDispatch()
  const classes = useStyles()
  const usersState = useSelector(state => state.users)

  const [loading, setLoading] = useState(false)

  const [filterUser, setFilterUser] = useState({})
  const [users, setUsers] = useState([])

  const handleRequestFilter = (property, value) => {
    const user = { ...filterUser }
    if (value === '') {
      delete user[property]
    } else {
      user[property] = value.toString()
    }
    setFilterUser(user)
  }

  useEffect(() => {
    setUsers(Object.values(usersState.entities)
      .filter(user => {
        for (const [key, value] of Object.entries(filterUser)) {
          if (!user[key].toString()
            .toLowerCase()
            .includes(value)) {
            return false
          }
        }
        return true
      }))
  }, [usersState, filterUser])

  useEffect(() => {
    if (loading) {
      if (usersState.loadedData) setLoading(false)
      return
    }

    if (usersState.loadedData) return

    setLoading(true)
    dispatch(getUsers(null))
  }, [usersState])

  const [openEditDialog, setOpenEditDialog] = useState(false)
  const [create, setCreate] = useState(false)
  const [userData, setUserData] = useState(null)

  const openEdit = (create, data) => {
    setCreate(create)
    setOpenEditDialog(true)
    if (create) return

    setUserData(data)
  }

  const [openPasswordDialog, setOpenPasswordDialog] = useState(false)
  const openPassword = data => {
    setUserData(data)
    setOpenPasswordDialog(true)
  }

  const [openDeleteDialog, setOpenDeleteDialog] = useState(false)
  const openDelete = data => {
    setUserData(data)
    setOpenDeleteDialog(true)
  }

  const [openAddDialog, setOpenAddDialog] = useState(false)
  const openAdd = data => {
    setUserData(data)
    setOpenAddDialog(true)
  }

  return (
    <div className={classes.div}>
      <EditUserDialog
        open={openEditDialog}
        onClose={() => setOpenEditDialog(false)}
        create={create}
        userData={userData}
      />
      <ChangePasswordDialog
        open={openPasswordDialog}
        onClose={() => setOpenPasswordDialog(false)}
        id={userData ? userData.id : 0}
      />
      <DeleteUserDialog
        open={openDeleteDialog}
        onClose={() => setOpenDeleteDialog(false)}
        id={userData ? userData.id : 0}
      />
      <AddUserRoleDialog
        open={openAddDialog}
        onClose={() => setOpenAddDialog(false)}
        userData={userData}
      />
      {
        loading
          ? <LinearProgress/>
          : ''
      }
      <CreateUserBar openEditDialog={openEdit}/>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              {
                columns.map(column => (
                  <TableCell
                    key={column.field}
                    align='center'
                  >
                    {
                      column.field === 'is_admin'
                        ? (
                          <FormControl>
                            <InputLabel>Admin</InputLabel>
                            <Select className={classes.select} onChange={e => handleRequestFilter(column.field, e.target.value)}>
                              <MenuItem value=''>all</MenuItem>
                              <MenuItem value='true'>true</MenuItem>
                              <MenuItem value='false'>false</MenuItem>
                            </Select>
                          </FormControl>)
                        : (
                          <TextField
                            label={'search by ' + column.headerName}
                            onChange={e => handleRequestFilter(column.field, e.target.value.toLowerCase())}
                          />)
                    }
                  </TableCell>
                ))
              }
            </TableRow>
          </TableHead>
        </Table>
      </TableContainer>
      {
        users.map(user =>
          <UserInfo
            key={user.id}
            userData={user}
            openEdit={openEdit}
            openPassword={openPassword}
            openDelete={openDelete}
            openAdd={openAdd}
          />)
      }
    </div>
  )
}
