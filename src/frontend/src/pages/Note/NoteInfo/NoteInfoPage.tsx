import { VStack, Text, Box, Switch } from "@chakra-ui/react"
import GetNote from "postAPI/notes/Get"
import React from "react"
import { NavigateFunction, Params } from "react-router-dom"
import { Note as NoteT} from "types/Note"
import styles from "./NotesInfoPage.module.scss";
import RoundButton from "components/RoundButton/RoundButton"

import CreateNote from "postAPI/tickets/Create"

type State = {
    note?: NoteT
    useBonusPoints: boolean
}

type NoteInfoParams = {
    match: Readonly<Params<string>>
    navigate: NavigateFunction
}

class NoteInfoPage extends React.Component<NoteInfoParams, State> {
    id: string;

    constructor(props) {
        super(props);
        this.id = this.props.match.id || "?";
        this.state = {
            useBonusPoints: false
        }
    };

    componentDidMount(): void {
        GetNote(this.id).then(data => {
            console.log(data)
            if (data.status === 200) {
                this.setState({note: data.content})
            }
        });
    }

    submit(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
        if (!this.state.note) {
            console.error("note is null")
            return
        }

        let button = e.currentTarget
        button.disabled = true
        CreateNote(this.state.note.id, "Заголовок № " + (this.state.note.id + 1), "Содержание заметки", "17.06.2023 13:17:04")
        .then(data => {
            button.disabled = false
            if (data.status === 200) {
                window.location.href = '/notes';
            } else {
                var title = document.getElementById("undertitle")
                if (title)
                    title.innerText = "Новая заметка создана"
            }
        });
    }

    async handleToggle() {
        this.setState({useBonusPoints: !this.state.useBonusPoints})
    }

    render() {
        return (
            <VStack className={styles.main_box}>
                {this.state.note &&
                    <Box style={{width: "100%"}}>
                        <Box className={styles.basics_box}>
                        <Text>Идентификатор: {this.id}</Text>
                        </Box>
                        <Box className={styles.basics_box}>
                        <Text>Заголовок: {this.state.note?.Title}</Text>
                        </Box>
                        <Box className={styles.basics_box}>
                        <Text>Содержание: {this.state.note?.Content}</Text>
                        </Box>
                        <Box className={styles.basics_box}>
                        <Text>Создано: {this.state.note?.CreatedAt}</Text>
                        </Box>
{/* 
                        <Box className={styles.bonus_div}>
                            <Switch
                                isChecked={this.state.useBonusPoints}
                                onChange={() => this.handleToggle()}
                                colorScheme="teal"
                                size="md"
                            />
                            <Text>
                                Использовать бонусные баллы
                            </Text>
                        </Box> */}


                        <RoundButton className={styles.basics_button} type="submit" onClick={event => this.submit(event)}>
                            Создать еще заметку
                        </RoundButton>
                    </Box>
                }
            </VStack>
        );
    };
};

export default NoteInfoPage;