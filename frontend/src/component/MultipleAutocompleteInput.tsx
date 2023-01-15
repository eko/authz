import * as React from 'react';
import { FilterOptionsState, InputAdornment, SxProps } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import CloseIcon from '@mui/icons-material/Close';
import SearchIcon from '@mui/icons-material/Search';
import Autocomplete, {
  AutocompleteCloseReason,
  createFilterOptions,
} from '@mui/material/Autocomplete';
import Box from '@mui/material/Box';
import { Button, List, ListItem, ListItemText, TextField } from '@mui/material';

import './MultipleAutocompleteInput.css';

export type ItemType = {
  id: string;
  label: string;
  raw?: any;
}

type FetcherFunc = (input: string) => Promise<ItemType[]>
type SetValueFunc = (items: ItemType[]) => void;

type AutocompleteInputProps = {
  disabled?: boolean
  allowAdd?: boolean
  defaultValues?: ItemType[]
  errorField?: any
  fetcher: FetcherFunc
  label?: string
  placeholder?: string
  setValue: SetValueFunc
  inputSx?: SxProps
  style?: React.CSSProperties
}

export default function MultipleAutocompleteInput({
  disabled,
  allowAdd,
  defaultValues,
  errorField,
  fetcher,
  label,
  placeholder,
  setValue,
  inputSx,
  style,
}: AutocompleteInputProps) {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const [listItems, setListItems] = React.useState<ItemType[]>([]);
  const [items, setItems] = React.useState<ItemType[]>([]);
  const theme = useTheme();

  const handleOnKeyUp = async (event: any) => {
    const values = await fetcher(event.target.value);
    setListItems(values);
  };

  const handleClose = () => {
    if (anchorEl) {
      anchorEl.focus();
    }
    setAnchorEl(null);
  };

  const handleOnClose = (event: React.ChangeEvent<{}>, reason: AutocompleteCloseReason) => {
    if (reason === 'escape') {
      handleClose();
    }
  };

  const handleOnChange = (event: React.SyntheticEvent, newItems: any, reason: string) => {
    if (event.type === 'keydown' &&
      (event as React.KeyboardEvent).key === 'Backspace' &&
      reason === 'removeOption') {
      return;
    }

    setItems(newItems);
  };

  React.useEffect(() => {
    setValue(items);
  }, [items, defaultValues, setValue]);

  React.useEffect(() => {
    if (defaultValues === undefined || defaultValues.length === 0) {
      return;
    }

    setListItems(defaultValues);
    setItems(defaultValues);
  }, [defaultValues]);

  const filter = createFilterOptions<ItemType>();

  const filterOptions = (options: ItemType[], params: FilterOptionsState<ItemType>) => {
    const filtered = filter(options, params);

    const { inputValue } = params;
    const isExisting = options.some((option) => inputValue === option.id);
    if (inputValue !== '' && !isExisting) {
      filtered.push({
        id: inputValue,
        label: inputValue,
      });
    }

    return filtered;
  };

  return (
    <div style={style}>
      <Autocomplete
        disabled={disabled}
        multiple
        openOnFocus
        freeSolo={allowAdd}
        selectOnFocus
        clearOnBlur
        onClose={handleOnClose}
        value={items}
        onChange={handleOnChange}
        disableCloseOnSelect
        filterOptions={allowAdd ? filterOptions : undefined}
        isOptionEqualToValue={(option: ItemType, value: ItemType): boolean => {
          return option.id === value.id;
        }}
        renderTags={() => null}
        noOptionsText='Aucun élément'
        renderOption={(props, option, { selected }) => (
          <li {...props}>
            <Box
              sx={{
                flexGrow: 1,
                '& span': {
                  color: theme.palette.mode === 'light' ? '#586069' : '#8b949e',
                },
              }}
            >
              {option.label}
            </Box>
            <Box
              component={CloseIcon}
              sx={{ opacity: 0.6, width: 18, height: 18 }}
              style={{ visibility: selected ? 'visible' : 'hidden' }}
            />
          </li>
        )}
        options={[...listItems]}
        getOptionLabel={(option) => typeof(option) === 'string' ? option : option.label}
        onFocus={handleOnKeyUp}
        onKeyUp={handleOnKeyUp}
        renderInput={(params) => (
          <TextField
            ref={params.InputProps.ref}
            disabled={disabled}
            role='multipleautocompleteinput-field'
            inputProps={params.inputProps}
            label={label}
            placeholder={placeholder}
            error={errorField ? true : false}
            helperText={errorField?.message}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
            sx={inputSx}
          />
        )}
      />

      {items.length > 0 ? (
        <List role='selected-list' sx={{ maxWidth: 360, bgcolor: 'background.paper' }}>
          {items.map((item, index) => (
            <ListItem role='multipleautocompleteinput-item' key={index} dense>
              <ListItemText primary={item.label} />

              {!disabled ? (
                <Button onClick={() => setItems(items.filter((_, i) => i !== index))} sx={{ right: -12 }}>
                  <Box
                    component={CloseIcon}
                    sx={{ opacity: 0.6, width: 18, height: 18 }}
                  />
                </Button>
              ) : null}
            </ListItem>
          ))}
        </List>
      ) : (
        <p style={{ margin: '10px 0 0 2px' }}>Aucun élément sélectionné.</p>
      )}
    </div>
  );
}
