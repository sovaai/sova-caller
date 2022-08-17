import { InputChangeEventHandlerType } from '../@types/common'

export const createInputOnChangeHandler =
  (setter: (value: any) => void): InputChangeEventHandlerType =>
  (e) =>
    setter(e.target.value)
