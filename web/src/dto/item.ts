// src/dto/item.ts
export type ItemType = 'folder' | 'file'

export interface Item {
    id: string
    name: string
    type: ItemType
    size?: number
    children?: Item[]
}
