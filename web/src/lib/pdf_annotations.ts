export type AnnotationTool = 'select' | 'textbox'

export type TextboxAnnotationData = {
  text: string
  left: number
  top: number
  width: number
  fontSize: number
  fill: string
  angle?: number
}

export type PDFAnnotation = {
  id: string
  type: 'textbox'
  data: TextboxAnnotationData
}

export const annotationTools: Array<{
  id: AnnotationTool
  label: string
  description: string
}> = [
  {
    id: 'select',
    label: 'Select',
    description: 'Move and edit existing annotations.',
  },
  {
    id: 'textbox',
    label: 'Text Box',
    description: 'Create a text annotation on the current page.',
  },
]

const defaultAnnotations: Record<number, PDFAnnotation[]> = {
  1: [
    {
      id: 'default-welcome',
      type: 'textbox',
      data: {
        text: 'Add notes here',
        left: 0.08,
        top: 0.08,
        width: 0.28,
        fontSize: 0.032,
        fill: '#0f172a',
      },
    },
  ],
}

export async function getAllAnnotations(page: number): Promise<PDFAnnotation[]> {
  const annotations = defaultAnnotations[page] ?? []
  return annotations.map(cloneAnnotation)
}

export async function updateAnnotation(page: number, data: PDFAnnotation): Promise<void> {
  void page
  void data
}

export function cloneAnnotation(annotation: PDFAnnotation): PDFAnnotation {
  return {
    ...annotation,
    data: { ...annotation.data },
  }
}
