const formatTel = (tel: string) => {
  const cleaned = tel.replace(/\D/g, '')

  const match = cleaned.match(/^(7|8)(\d{3})(\d{3})(\d{2})(\d{2})$/)

  if (match) {
    return ['+7', match[2], match[3], match[4], match[5]].join('')
  }

  return ''
}

export const parseContacts = (rawContacts: string): string[] => {
  const contacts: string[] = rawContacts
    .split(';')
    .map((r) => formatTel(r))
    .filter((tel) => tel.length !== 0)

  return contacts
}
