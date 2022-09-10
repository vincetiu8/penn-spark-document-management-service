export const toggleLoad = state => {
  state.loading = true
  state.error = undefined
}

export const onError = (state, action) => {
  state.loading = false
  state.error = action.payload
}
