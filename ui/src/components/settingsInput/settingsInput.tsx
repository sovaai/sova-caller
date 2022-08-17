import { FormControl, Input } from '@chakra-ui/react'
import React, { memo } from 'react'

import { createInputOnChangeHandler } from '../../helpers/createInputOnChangeHandler'
import { colors } from '../../styles/theme'

export interface SettingsField {
  value: string
  setter: (value: any) => void
  placeholder: string
}

export interface SettingsInputProps extends SettingsField {
  disabled: boolean
}

const SettingsInput = memo(({ value, setter, placeholder, disabled }: SettingsInputProps) => {
  return (
    <FormControl>
      <Input
        value={value}
        onChange={createInputOnChangeHandler(setter)}
        placeholder={placeholder}
        size="md"
        borderRadius={'2.5rem'}
        borderColor={colors[4]}
        backgroundColor={colors[4]}
        fontSize={'14px'}
        lineHeight={'20px'}
        fontWeight={400}
        color={colors[2]}
        _placeholder={{
          color: colors[3],
        }}
        spellCheck={false}
        disabled={disabled}
      />
    </FormControl>
  )
})

SettingsInput.displayName = 'SettingsInput'

export { SettingsInput }
