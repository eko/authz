import { useState } from 'react';
import {
    Button, Dialog, DialogActions,
    DialogContent, DialogContentText, DialogTitle,
} from '@mui/material';

const useConfirm = () => {
  const [title, setTitle] = useState('Confirmation');
  const [message, setMessage] = useState('Êtes-vous certain de vouloir effectuer cette action ?');

  const [promise, setPromise] = useState<{ resolve: Function, reject?: Function } | null>(null);
    
    const confirm = (title: string, message: string) => new Promise((resolve: Function, reject: Function) => {
        setTitle(title);
        setMessage(message);
        setPromise({ resolve });
    });
    
    const handleClose = () => {
        setPromise(null);
    };
    
    const handleConfirm = () => {
        promise?.resolve(true);
        handleClose();
    };
    
    const handleCancel = () => {
        promise?.resolve(false);
        handleClose();
    };

    const ConfirmationDialog = () => (
      <Dialog
        open={promise !== null}
        fullWidth
      >
        <DialogTitle>{title}</DialogTitle>
        <DialogContent>
          <DialogContentText>{message}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button variant='contained' onClick={handleConfirm}>OK</Button>
          <Button onClick={handleCancel}>Annuler</Button>
        </DialogActions>
      </Dialog>
    );

    return {
      ConfirmationDialog,
      confirm
    };
};
  
  export default useConfirm;