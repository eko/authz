import { useState } from 'react';
import { IconButton, Chip } from '@mui/material';
import CheckCircleOutlineOutlinedIcon from '@mui/icons-material/CheckCircleOutlineOutlined';
import ContentCopyOutlinedIcon from '@mui/icons-material/ContentCopyOutlined';

type CopyToClipboardButtonProps = {
    text: string
    color?: 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning'
    variant?: 'filled' | 'outlined'
}

export default function CopyToClipboardButton({
    text,
    color = 'primary',
    variant = 'outlined',
}: CopyToClipboardButtonProps) {
    const [copied, setCopied] = useState(false);

    const handleOnClick = () => {
      navigator.clipboard.writeText(text);

      setCopied(true);
      setTimeout(() => setCopied(false), 1000);
    }
    
    return (
        <>
          <Chip label={text} size='small' variant={variant} color={color} sx={{ mr: 0.5 }} />

          <IconButton
            role='copy-icon'
            size='small'
            color={color}
            aria-label='Copier dans le presse-papier'
            onClick={handleOnClick}
          >
            <ContentCopyOutlinedIcon sx={{ fontSize: '1.0rem' }} />
          </IconButton>

          {copied ? (
            <Chip
              role='copy-confirmation'
              icon={<CheckCircleOutlineOutlinedIcon />}
              size='small'
              label='Copié'
              sx={{ ml: 0.5 }} />
          ) : null}
        </>
    )
};