import { Flex, Heading, Image, Text, useStyleConfig } from '@chakra-ui/react'
import React, { memo } from 'react'

import { colors } from '../../styles/theme'
import logo from '../../assets/logo.svg'

const Header = memo(() => {
  const headerStyles = useStyleConfig('Header')

  return (
    <Flex as={'header'} sx={headerStyles}>
      <Image src={logo} alt="logo" />
      <Heading
        as="h1"
        fontFamily={`'Montserrat', sans-serif`}
        fontSize={'1.875rem'}
        fontWeight={'600'}
        lineHeight={'1.5rem'}
      >
        <Text color={colors[2]} as="span" margin={2}>
          SOVA
        </Text>
        <Text color={colors[1]} as="span">
          Caller
        </Text>
      </Heading>
    </Flex>
  )
})

Header.displayName = 'Header'

export { Header }
