import React from "react";
import { RequestStat } from "types/Statistics";

interface Props {
  requestStats: RequestStat[];
}

const PopularNotes: React.FC<Props> = ({ requestStats }) => {
  const countRequestsPerNote = (requestStats: RequestStat[]): Record<string, number> => {
    const noteCounts: Record<string, number> = {};

    requestStats.forEach((stat) => {
      const match = stat.path.match(/\/api\/v1\/notes\/(\w+)/);
      if (match) {
        const id = match[1];
        if (noteCounts[id]) {
          noteCounts[id]++;
        } else {
          noteCounts[id] = 1;
        }
      }
    });

    return noteCounts;
  };

  const noteCounts = countRequestsPerNote(requestStats);

  // Получаем массив в формате [{ id: string, count: number }]
  const notes = Object.entries(noteCounts).map(([id, count]) => ({
    id,
    count,
  }));

  // Сортируем рейсы по убыванию количества обращений
  const sortedNotes = notes.sort((a, b) => b.count - a.count);

  const topNotes = sortedNotes.slice(0, 3);

  return (
    <div>
      <h2>Самые важные заметки</h2>
      <ul>
        {topNotes.map((note) => (
          <li key={note.id}>
            Заметка № {note.id}: {note.count} просмотров
          </li>
        ))}
      </ul>
    </div>
  );
};

export default PopularNotes;
