import { Flex, Heading, Textarea } from '@chakra-ui/react'
import React, { memo } from 'react'

import { createTextareaOnChangeHandler } from '../../helpers/createTextareaOnChangeHandler'
import { colors } from '../../styles/theme'

interface ContactsProps {
  value: string
  setter: (value: any) => void
  disabled: boolean
}

const Contacts = memo(({ value, setter, disabled }: ContactsProps) => {
  return (
    <Flex w="100%" flexDirection={'column'} height={'100%'} justifyContent={'flex-start'}>
      <Heading
        as="h2"
        size="xs"
        w="100%"
        pl={'1.25rem'}
        fontSize={'14px'}
        lineHeight={'20px'}
        fontWeight={600}
        color={colors[2]}
      >
        Список телефонов
      </Heading>
      <Textarea
        value={value}
        onChange={createTextareaOnChangeHandler(setter)}
        placeholder="Вставьте или запишите здесь все номера через “;”"
        resize={'none'}
        height={'100%'}
        mt={2}
        px="1.5rem"
        py="1.25rem"
        backgroundColor={colors[4]}
        borderColor={colors[4]}
        borderRadius={'20px'}
        fontSize={'14px'}
        lineHeight={'1.25rem'}
        fontWeight={400}
        color={colors[2]}
        _placeholder={{
          color: colors[3],
        }}
        spellCheck={false}
        disabled={disabled}
      />
    </Flex>
  )
})

Contacts.displayName = 'Contacts'

export { Contacts }
