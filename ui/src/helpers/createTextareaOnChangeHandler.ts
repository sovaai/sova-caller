import { TextareaChangeEventHandlerType } from '../@types/common'

export const createTextareaOnChangeHandler =
  (setter: (value: any) => void): TextareaChangeEventHandlerType =>
  (e) =>
    setter(e.target.value)
