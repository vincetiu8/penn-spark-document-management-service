import React from 'react'
import {
  Button,
  ButtonGroup,
  Card,
  CardActions,
  CardContent,
  Grid,
  Typography
} from '@material-ui/core'
import { makeStyles } from '@material-ui/core/styles'
import AssignmentIndIcon from '@material-ui/icons/AssignmentInd'
import { useDispatch } from 'react-redux'
import { removeUserRole } from '../../store/usersSlice'

const useStyles = makeStyles({
  card: {
    verticalAlign: 'middle',
    display: 'inline-flex',
    width: '100%'
  },
  cardActionArea: {
    paddingLeft: '2ch',
    display: 'flex',
    flexGrow: 1,
    width: 'auto'
  },
  cardActions: {
    float: 'left'
  },
  button: {
    textTransform: 'none',
    fontSize: '1.25rem',
    width: '100%',
    whiteSpace: 'nowrap'
  },
  typography: {
    fontSize: '1.25rem'
  }
})

export const UserRoleInfo = ({
  userRoleData,
  userData
}) => {
  const classes = useStyles()
  const dispatch = useDispatch()

  const onClick = () => {
    dispatch(removeUserRole({
      userID: userData.id,
      userRole: userRoleData
    }))
  }

  return (
    <Grid item>
      <Card variant='outlined' square className={classes.card}>
        <CardContent
          className={classes.cardActionArea}
        >
          <Grid container spacing={2} className={classes.grid}>
            <Grid item>
              <AssignmentIndIcon/>
            </Grid>
            <Grid item>
              <Typography className={classes.typography}>
                {userRoleData.name}
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
        <CardActions className={classes.cardActions}>
          <ButtonGroup>
            <Button className={classes.button} onClick={onClick}>
              Remove User Role
            </Button>
          </ButtonGroup>
        </CardActions>
      </Card>
    </Grid>
  )
}
