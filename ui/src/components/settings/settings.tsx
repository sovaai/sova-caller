import { Heading, VStack } from '@chakra-ui/react'
import React, { memo } from 'react'

import { colors } from '../../styles/theme'
import { SettingsField, SettingsInput } from '../settingsInput/settingsInput'

interface SettingsProps {
  heading: string
  fields: SettingsField[]
  disabled: boolean
}

const Settings = memo(({ heading, fields, disabled }: SettingsProps) => {
  return (
    <VStack spacing={1.5} width={'100%'}>
      <Heading
        as="h2"
        size="xs"
        w="100%"
        pl={'1.25rem'}
        fontSize={'14px'}
        lineHeight={'20px'}
        fontWeight={600}
        color={colors[2]}
        pb={0.5}
      >
        {heading}
      </Heading>
      {fields.map((field, index) => (
        <SettingsInput
          key={index}
          value={field.value}
          setter={field.setter}
          placeholder={field.placeholder}
          disabled={disabled}
        />
      ))}
    </VStack>
  )
})

Settings.displayName = 'Settings'

export { Settings }
