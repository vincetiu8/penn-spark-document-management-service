import React, { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { getUserRoles } from '../../store/userRolesSlice'
import { UserRolesBar } from './UserRolesBar'
import { UserRoleInfo } from './UserRoleInfo'
import { EditUserRoleDialog } from './EditUserRoleDialog'
import { DeleteUserRoleDialog } from './DeleteUserRoleDialog'
import { AddAccessRoleDialog } from './AddAccessRoleDialog'

export const UserRoles = () => {
  const dispatch = useDispatch()
  const userRolesState = useSelector(state => state.userRoles)

  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (userRolesState.loading) return

    if (loading) {
      setLoading(false)
      return
    }

    if (userRolesState.loadedData) return

    setLoading(true)
    dispatch(getUserRoles(null))
  }, [userRolesState])

  const [searchTerm, setSearchTerm] = useState('')
  const [userRoles, setUserRoles] = useState([])

  useEffect(() => {
    setUserRoles(Object.values(userRolesState.entities)
      .filter(userRole => userRole.name.toLowerCase()
        .includes(searchTerm)))
  }, [searchTerm, userRolesState])

  const [create, setCreate] = useState(false)
  const [userRoleData, setUserRoleData] = useState(null)
  const [openEditDialog, setOpenEditDialog] = useState(false)

  const openEdit = (create, userRoleData) => {
    setCreate(create)
    setUserRoleData(userRoleData)
    setOpenEditDialog(true)
  }

  const [openDeleteDialog, setOpenDeleteDialog] = useState(false)
  const [deleteID, setDeleteID] = useState(null)

  const openDelete = id => {
    setOpenDeleteDialog(true)
    setDeleteID(id)
  }

  const [openAddDialog, setOpenAddDialog] = useState(false)
  const openAdd = data => {
    setUserRoleData(data)
    setOpenAddDialog(true)
  }

  return (
    <div>
      <EditUserRoleDialog open={openEditDialog} setOpen={setOpenEditDialog} create={create}
                          userRoleData={userRoleData}/>
      <DeleteUserRoleDialog open={openDeleteDialog} onClose={() => setOpenDeleteDialog(false)}
                            id={deleteID}/>
      <AddAccessRoleDialog open={openAddDialog} onClose={() => setOpenAddDialog(false)}
                           userRoleData={userRoleData}/>
      <UserRolesBar searchTerm={searchTerm} setSearchTerm={setSearchTerm}
                    openEdit={openEdit}/>
      {
        userRoles.map(userRole => <UserRoleInfo
          key={userRole.id}
          userRoleData={userRole}
          openEdit={openEdit}
          openDelete={openDelete}
          openAdd={openAdd}
        />)
      }
    </div>
  )
}
